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

func GetMultiVideo(latestTime time.Time, count int) []dao.Video {
	videoInfoList := dao.QueryVideosByPublishTime(latestTime, count)
	return videoInfoList
}

func StoreVideo(data []byte, username string, fileName string) error {
	// 暂时存在本地，TODO:对象存储
	fileName = time.Now().Truncate(time.Second).Format("20060102150405") +
		"_" + username + fileName

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

func QueryVideoById(videoID uint64) dao.Video {
	video, _ := dao.QueryVideoByID(videoID)
	return video
}
