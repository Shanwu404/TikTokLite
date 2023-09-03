package service

import (
	"fmt"
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"testing"
	"time"
)

func CommentServiceImplInit() {
	dao.Init()
	redis.InitRedis()
}

func TestCommentServiceImpl_QueryCommentsByVideoId(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	comments := csi.QueryCommentsByVideoId(8)
	fmt.Println(comments)
}

func TestCommentServiceImpl_PostComment(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	comment := CommentParams{
		UserId:     4000,
		VideoId:    4000,
		Content:    "test",
		CreateDate: time.Now(),
	}
	id, code, messgae := csi.PostComment(comment)
	fmt.Println(id, code, messgae)
}

func TestCommentServiceImpl_DeleteComment(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	code, message := csi.DeleteComment(59)
	fmt.Println(code, message)
}

func TestCommentServiceImpl_CountComments(t *testing.T) {
	CommentServiceImplInit()
	csi := CommentServiceImpl{}
	cnt := csi.CountComments(0)
	fmt.Println(cnt)
}
