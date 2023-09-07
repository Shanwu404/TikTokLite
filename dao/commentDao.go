package dao

import (
	"time"

	"github.com/Shanwu404/TikTokLite/log/logger"
)

type Comment struct {
	Id         int64
	UserId     int64
	VideoId    int64
	Content    string
	CreateDate time.Time
}

// QueryCommentsByVideoId 根据视频id查询评论列表
func QueryCommentsByVideoId(id int64) ([]Comment, error) {
	var comments []Comment
	if err := db.Where("video_id = ?", id).Order("create_date DESC").Find(&comments).Error; err != nil {
		logger.Errorln(err.Error())
		return comments, err
	}
	return comments, nil
}

// InsertComment 插入评论
func InsertComment(comment Comment) (Comment, error) {
	if err := db.Model(Comment{}).Create(&comment).Error; err != nil {
		logger.Errorln(err.Error())
		return Comment{}, err
	}
	return comment, nil
}

// DeleteComment 根据评论id删除评论
func DeleteComment(id int64) bool {
	var comment Comment
	if err := db.Where("id = ?", id).First(&comment).Error; err != nil {
		logger.Errorln(err.Error())
		return false
	}
	db.Delete(&comment)
	return true
}

// CountComments 根据视频id统计评论数量
func CountComments(id int64) (int64, error) {
	var count int64
	err := db.Model(&Comment{}).Where("video_id = ?", id).Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}
