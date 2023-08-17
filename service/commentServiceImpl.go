package service

import (
	"log"

	"github.com/Shanwu404/TikTokLite/dao"
)

type CommentServiceImpl struct{}

func NewCommentService() CommentService {
	return &CommentServiceImpl{}
}

func (CommentServiceImpl) QueryCommentsByVideoId(id int64) []dao.Comment {
	comments, err := dao.QueryCommentsByVideoId(id)
	if err != nil {
		log.Println("error:", err.Error())
		return comments
	}
	log.Println("Query comments successfully!")
	return comments
}

func (CommentServiceImpl) PostComment(comment dao.Comment) (int64, int32, string) {
	comment, err := dao.InsertComment(comment)
	if err != nil {
		return -1, 1, "Post comment failed!"
	}
	return comment.Id, 0, "Post comment successfully!"
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
		log.Println("count from db error:", err)
		return 0
	}
	log.Println("count comments successfully!")
	return cnt
}
