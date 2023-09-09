package dao

import (
	"time"

	"github.com/Shanwu404/TikTokLite/log/logger"
)

type Video struct {
	ID          int64
	AuthorID    int64
	PlayURL     string
	CoverURL    string
	PublishTime time.Time
	Title       string
}

func InsertVideo(video Video) error {
	err := db.Create(&video).Error
	if err != nil {
		logger.Errorln(err)
	}
	return err
}

func QueryVideoByID(id int64) (Video, error) {
	video := Video{}
	err := db.Take(&video, id).Error // Take查不到会返回err
	if err != nil {
		logger.Errorln(err, id)
	}
	return video, err
}

func QueryVideosByPublishTime(latestTime time.Time, count int) []Video {
	videos := make([]Video, 0, count)
	if err := db.
		Where("publish_time < ?", latestTime).
		Order("publish_time desc").
		Limit(count).Find(&videos).Error; err != nil {
		logger.Errorln(err)
	}
	return videos
}

func QueryVideosByAuthorId(authorId int64) []Video {
	var videos = make([]Video, 0)
	if err := db.Where("author_id = ?", authorId).Find(&videos).Error; err != nil {
		logger.Errorln(err, authorId)
	}
	return videos
}
