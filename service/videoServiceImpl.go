package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils/aliyun/ossClient"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	internal = "Internal"
	external = "External"
)

// var videoMutex sync.RWMutex

type VideoServiceImpl struct {
}

func NewVideoService() VideoService {
	return &VideoServiceImpl{}
}

func (vService *VideoServiceImpl) Exist(videoID int64) bool {
	redisKey := "video:" + strconv.FormatInt(videoID, 10)
	exist := redis.RDb.Exists(redis.Ctx, redisKey).Val() == 1
	if !exist {
		video, err := dao.QueryVideoByID(videoID)
		if err != nil {
			return false
		}
		exBucket, err := ossClient.NewBucket(external)
		if err != nil {
			logger.Errorln("Get bucket error:", err)
			return true
		}
		externalURL, err := exBucket.ObjectExternalURL(video.PlayURL)
		if err != nil {
			logger.Errorln("Get video play url error:", err)
			return true
		}
		redis.RDb.SetNX(redis.Ctx, redisKey, externalURL, time.Second*(ossClient.SignedURLExpiration-60))
	}
	return true
}

func (vService *VideoServiceImpl) GetMultiVideoBefore(latestTimestamp int64, count int) []VideoParams {
	mb, err := ossClient.NewBucket(external)
	if err != nil {
		logger.Errorln("Get bucket error:", err)
		return []VideoParams{}
	}
	latestTime := time.Unix(latestTimestamp, 0)
	videos := dao.QueryVideosByPublishTime(latestTime, count)
	videoSlc := make([]VideoParams, 0, len(videos))
	for i := range videos {
		externalPlayURL, ok := getPlayURLFromRedis(videos[i].ID)
		if ok {
			videos[i].PlayURL = externalPlayURL
		} else {
			videos[i].PlayURL, err = mb.ObjectExternalURL(videos[i].PlayURL)
			if err != nil {
				logger.Errorln("Get object url error:", err)
			}
		}
		videos[i].CoverURL, err = mb.ObjectExternalURL(videos[i].CoverURL)
		if err != nil {
			logger.Errorln("Get video cover url error:", err)
		}
		videoSlc = append(videoSlc, VideoParams(videos[i]))
	}
	return videoSlc
}

func (vService *VideoServiceImpl) StoreVideo(dataHeader *multipart.FileHeader, fileName string, video *VideoParams) error {
	mb, err := ossClient.NewBucket(internal)
	if err != nil {
		logger.Errorln("Get bucket error:", err)
		return err
	}
	internalURL := "videos/" + fileName
	err = mb.UploadVideo(dataHeader, internalURL)
	if err != nil {
		logger.Errorln("Upload video failed:", err)
		return err
	}
	video.PlayURL = internalURL

	frame := captureFrame(dataHeader, fileName)
	if frame != nil {
		internalImageURL := "images/" + fileName
		err = mb.PutObject(internalImageURL, frame)
		if err != nil {
			logger.Errorln("Upload frame failed:", err)
		}
		video.CoverURL = internalImageURL
	}

	err = vService.InsertVideosTable(video)
	if err != nil {
		logger.Errorln("Insert video table failed:", err)
		return err
	}
	return nil
}

func (vService *VideoServiceImpl) InsertVideosTable(video *VideoParams) error {
	retry := 5
	delay := time.Second
	var err error
	for i := 0; i < retry; i++ {
		err = dao.InsertVideo(dao.Video(*video))
		if err != nil {
			time.Sleep(delay)
		} else {
			return nil
		}
	}
	return err
}

// TODO: 分页查询
func (vService *VideoServiceImpl) GetVideoListByUserId(AuthorID int64) []VideoParams {
	mb, err := ossClient.NewBucket(external)
	if err != nil {
		logger.Errorln("Get bucket error:", err)
		return []VideoParams{}
	}
	videos := dao.QueryVideosByAuthorId(AuthorID)
	results := make([]VideoParams, 0, len(videos))
	for i := range videos {
		externalPlayURL, ok := getPlayURLFromRedis(videos[i].ID)
		if ok {
			videos[i].PlayURL = externalPlayURL
		} else {
			videos[i].PlayURL, err = mb.ObjectExternalURL(videos[i].PlayURL)
			if err != nil {
				logger.Errorln("Get object url error:", err)
			}
		}
		videos[i].CoverURL, err = mb.ObjectExternalURL(videos[i].CoverURL)
		if err != nil {
			logger.Errorln("Get video cover url error:", err)
		}
		results = append(results, VideoParams(videos[i]))
	}
	return results
}

func (vService *VideoServiceImpl) QueryVideoById(videoID int64) VideoParams {
	mb, err := ossClient.NewBucket(external)
	if err != nil {
		logger.Errorln("Get bucket error:", err)
		return VideoParams{}
	}
	videoFromDao, _ := dao.QueryVideoByID(videoID)
	// 视频实际链接很适合放入redis，设置过期时间
	externalPlayURL, ok := getPlayURLFromRedis(videoFromDao.ID)
	if ok {
		videoFromDao.PlayURL = externalPlayURL
	} else {
		videoFromDao.PlayURL, err = mb.ObjectExternalURL(videoFromDao.PlayURL)
		if err != nil {
			logger.Errorln("Get object url error:", err)
		}
	}
	videoFromDao.PlayURL, err = mb.ObjectExternalURL(videoFromDao.PlayURL)
	if err != nil {
		logger.Errorln("Get video play url error:", err)
	}
	videoFromDao.CoverURL, err = mb.ObjectExternalURL(videoFromDao.CoverURL)
	if err != nil {
		logger.Errorln("Get video cover url error:", err)
	}
	video := VideoParams(videoFromDao)
	return video
}

// ----------------------------------private---------------------------------------------

func captureFrame(dataHeader *multipart.FileHeader, filename string) io.Reader {
	tempVideoPath := "videos/" + filename
	tempVideo, err := os.Create(tempVideoPath)
	if err != nil {
		logger.Errorln(errors.Join(errors.New("保存临时本地视频文件失败：创建失败："), err))
		return nil
	}
	defer os.Remove(tempVideoPath)

	videodata, err := dataHeader.Open()
	if err != nil {
		logger.Errorln(errors.Join(errors.New("保存临时本地视频文件失败：提取失败："), err))
		return nil
	}

	_, err = io.Copy(tempVideo, videodata)
	if err != nil {
		logger.Errorln(errors.Join(errors.New("保存临时本地视频文件失败：保存失败："), err))
		return nil
	}

	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(tempVideoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		logger.Errorln(errors.Join(errors.New("保存临时本地视频文件失败：创建失败："), err))
		return nil
	}
	return buf
}

func getPlayURLFromRedis(videoID int64) (string, bool) {
	redisKey := "video:" + strconv.FormatInt(videoID, 10)
	playURL, err := redis.RDb.Get(redis.Ctx, redisKey).Result()
	if err != nil {
		logger.Debugf("No PlayURL of videoID:%v in redis.\n", videoID)
		return "", false
	}
	logger.Debugf("The PlayURL of videoID:%v in redis is %v.\n", videoID, playURL)
	return playURL, true
}
