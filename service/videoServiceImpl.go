package service

import (
	"log"
	"mime/multipart"
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/utils/aliyun/ossClient"
)

type VideoServiceImpl struct {
}

func NewVideoService() VideoService {
	return &VideoServiceImpl{}
}

func (vService *VideoServiceImpl) GetMultiVideoBefore(latestTimestamp int64, count int) []VideoParams {
	latestTime := time.Unix(latestTimestamp, 0)
	videos := dao.QueryVideosByPublishTime(latestTime, count)
	videoSlc := make([]VideoParams, 0, len(videos))
	for i := range videos {
		videoSlc = append(videoSlc, VideoParams(videos[i]))
	}
	return videoSlc
}

func (vService *VideoServiceImpl) StoreVideo(dataHeader *multipart.FileHeader, fileName string, video *VideoParams) error {
	videoURL, err := ossClient.UploadVideo(dataHeader, fileName)
	if err != nil {
		log.Println("Upload video failed:", err)
		return err
	}
	video.PlayURL = videoURL
	err = vService.InsertVideosTable(video)
	if err != nil {
		log.Println("Insert video table failed:", err)
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
	videos := dao.QueryVideosByAuthorId(AuthorID)
	results := make([]VideoParams, 0, len(videos))
	for i := range videos {
		results = append(results, VideoParams(videos[i]))
	}
	return results
}

func (vService *VideoServiceImpl) QueryVideoById(videoID int64) VideoParams {
	videoFromDao, _ := dao.QueryVideoByID(videoID)
	video := VideoParams(videoFromDao)
	return video
}
