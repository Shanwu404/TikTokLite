package service

import (
	"mime/multipart"
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/utils/aliyun/ossClient"
)

const (
	internal = "Internal"
	external = "External"
)

type VideoServiceImpl struct {
}

func NewVideoService() VideoService {
	return &VideoServiceImpl{}
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
		videos[i].PlayURL, err = mb.ObjectExternalURL(videos[i].PlayURL)
		if err != nil {
			logger.Errorln("Get object url error:", err)
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
	err = vService.InsertVideosTable(video)
	if err != nil {
		logger.Errorln("Insert video table failed:", err)
		// 可优化的点：删除OSS视频，保持事务原子性
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
		videos[i].PlayURL, err = mb.ObjectExternalURL(videos[i].PlayURL)
		if err != nil {
			logger.Errorln("Get object url error:", err)
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
	videoFromDao.PlayURL, err = mb.ObjectExternalURL(videoFromDao.PlayURL)
	if err != nil {
		logger.Errorln("Get object url error:", err)
	}
	video := VideoParams(videoFromDao)
	return video
}
