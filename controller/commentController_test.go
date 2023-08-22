package controller

import (
	"strings"
	"testing"

	"github.com/Shanwu404/TikTokLite/middleware/auth"
)

func TestCommentAction(t *testing.T) {
	tok, _ := auth.GenerateToken("Kite", 1)
	token := "token=" + tok
	// 评论操作 - 添加评论
	url1 := "http://127.0.0.1:8080/douyin/comment/action/?video_id=1&action_type=1&comment_text=test"
	method1 := "POST"
	SendRequest(method1, url1, strings.NewReader(token))

	//评论操作 - 删除评论
	url2 := "http://127.0.0.1:8080/douyin/comment/action/?action_type=2&comment_id=27"
	method2 := "POST"
	SendRequest(method2, url2, strings.NewReader(token))
}

func TestCommentList(t *testing.T) {
	token, _ := auth.GenerateToken("Kite", 1)
	url1 := "http://127.0.0.1:8080/douyin/comment/list/?video_id=1&token=" + token
	method1 := "GET"
	SendRequest(method1, url1, strings.NewReader(token))
}
