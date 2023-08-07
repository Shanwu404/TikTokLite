package dao

import (
	"log"
	"time"
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
	if err := Db.Where("video_id = ?", id).Find(&comments).Error; err != nil {
		log.Println(err.Error())
		return comments, err
	}
	return comments, nil
}

// InsertComment 插入评论
func InsertComment(comment Comment) (Comment, error) {
	if err := Db.Model(Comment{}).Create(&comment).Error; err != nil {
		log.Println(err.Error())
		return Comment{}, err
	}
	return comment, nil
}

// DeleteComment 根据评论id删除评论
func DeleteComment(id int64) bool {
	var comment Comment
	if err := Db.Where("id = ?", id).First(&comment).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	Db.Delete(&comment)
	return true
}
