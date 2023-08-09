package dao

import (
	"log"
	"time"
)

type Video struct {
	ID          uint64
	AuthorID    uint64
	PlayURL     uint64
	CoverURL    uint64
	PublishTime time.Time
	Title       string
}

func InsertVideo(video Video) error {
	if err := db.Create(&video).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func QueryVideoByID(id uint64) (Video, error) {
	video := Video{}
	if err := db.First(&video, id).Error; err != nil {
		log.Println(err.Error())
		return video, err
	}
	return video, nil
}

func QueryVideosByPublishTime(latestTime time.Time, count int) []Video {
	videos := make([]Video, 0, count)
	if err := db.Order("publish_time desc").
		Where("publish_time < ?", latestTime).
		Limit(count).Find(&videos).Error; err != nil {
		log.Println(err.Error())
	}
	return videos
}

func QueryVideosByAuthorId(authorId uint64) []Video {
	var videos = make([]Video, 0)
	if err := db.Where("author_id = ?", authorId).Find(&videos).Error; err != nil {
		log.Println(err.Error())
	}
	return videos
}
