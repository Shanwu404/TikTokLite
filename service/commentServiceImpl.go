package service

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
)

type CommentServiceImpl struct{}

func NewCommentService() CommentService {
	return &CommentServiceImpl{}
}

func (CommentServiceImpl) QueryCommentsByVideoId(id int64) []CommentParams {
	comments, err := dao.QueryCommentsByVideoId(id)
	if err != nil {
		logger.Errorln("error:", err.Error())
	}
	results := make([]CommentParams, 0, len(comments))
	for i := range comments {
		results = append(results, CommentParams(comments[i]))
	}
	logger.Infoln("Query comments successfully!")
	return results
}

func (CommentServiceImpl) PostComment(comment CommentParams) (int64, int32, string) {
	commentNew, err := dao.InsertComment(dao.Comment(comment))
	if err != nil {
		return -1, 1, "Post comment failed!"
	}
	return commentNew.Id, 0, "Post comment successfully!"
}

func (CommentServiceImpl) DeleteComment(id int64) (int32, string) {
	flag := dao.DeleteComment(id)
	if flag == false {
		return 1, "Delete comment failed!"
	}
	return 0, "Delete comment successfully!"
}

func (CommentServiceImpl) CountComments(id int64) int64 {
	cnt, err := dao.CountComments(id)
	if err != nil {
		logger.Errorln("count from db error:", err)
		return 0
	}
	logger.Infoln("count comments successfully!")
	return cnt
}
