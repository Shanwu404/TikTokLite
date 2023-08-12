package service

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
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

func (vService *VideoServiceImpl) StoreVideo(data []byte, username string, fileName string) error {
	// 暂时存在本地，TODO:对象存储
	fileName = time.Now().Truncate(time.Second).Format("20060102150405") +
		"_" + username + "_" + fileName

	nameLimit := 999
	for i := 0; i <= nameLimit; i++ {
		if i > 0 {
			fileName = fileName + "(" + strconv.Itoa(i) + ")"
		}
		_, err := os.Stat(fileName)
		switch {
		case err == nil && i == nameLimit:
			log.Printf("File '%s' exists\n", fileName)
			return errors.New("too many files with that name already exists")
		case err == nil && i < nameLimit:
			continue
		case os.IsNotExist(err):
			log.Printf("Saving file as %v.\n", fileName)
			file, err := os.Create(fileName)
			if err != nil {
				log.Println("Error creating file:", err.Error())
				return err
			}
			defer file.Close()

			bytesWritten, err := io.Copy(file, bytes.NewReader(data))
			if err != nil {
				log.Println("Error writing data to file:", err.Error())
				return err
			}
			log.Println("Succeed! Bytes written:", bytesWritten)
			return nil
		default:
			log.Println("Error:", err.Error())
			return err
		}
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
