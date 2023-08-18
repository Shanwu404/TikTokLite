package service

import (
	"time"
)

type CommentParams struct {
	Id         int64
	UserId     int64
	VideoId    int64
	Content    string
	CreateDate time.Time
}

type CommentService interface {
	// QueryCommentsByVideoId 根据视频id查询评论列表
	QueryCommentsByVideoId(id int64) []CommentParams

	// PostComment 发布评论
	PostComment(comment CommentParams) (int64, int32, string)

	// DeleteComment 删除评论
	DeleteComment(id int64) (int32, string)

	// CountComments 统计评论数
	CountComments(id int64) int64
}
