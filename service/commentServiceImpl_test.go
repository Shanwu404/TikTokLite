package service

import (
	"fmt"
	"github.com/Shanwu404/TikTokLite/dao"
	"testing"
	"time"
)

func CommentServiceImplInit() {
	dao.Init()
}

func TestCommentServiceImpl_QueryCommentsByVideoId(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	comments := csi.QueryCommentsByVideoId(1000)
	fmt.Println(comments)
}

func TestCommentServiceImpl_PostComment(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	comment := dao.Comment{
		UserId:     1000,
		VideoId:    1000,
		Content:    "test",
		CreateDate: time.Now(),
	}
	id, code, messgae := csi.PostComment(comment)
	fmt.Println(id, code, messgae)
}

func TestCommentServiceImpl_DeleteComment(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	code, message := csi.DeleteComment(15)
	fmt.Println(code, message)
}
