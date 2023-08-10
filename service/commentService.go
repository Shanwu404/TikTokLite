package service

import "github.com/Shanwu404/TikTokLite/dao"

type CommentService interface {
	// QueryCommentsByVideoId 根据视频id查询评论列表
	QueryCommentsByVideoId(id int64) []dao.Comment

	// PostComment 发布评论
	PostComment(comment dao.Comment) (int64, int32, string)

	// DeleteComment 删除评论
	DeleteComment(id int64) (int32, string)
}
