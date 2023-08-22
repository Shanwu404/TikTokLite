package controller

import (
	"strings"
	"testing"

	"github.com/Shanwu404/TikTokLite/middleware/auth"
)

func TestFavoriteAction(t *testing.T) {
	tok, _ := auth.GenerateToken("chy", 7777)
	token := "token=" + tok
	url1 := "http://127.0.0.1:8080/douyin/favorite/action/?video_id=1&action_type=1"
	method1 := "POST"
	SendRequest(method1, url1, strings.NewReader(token))
}

func TestFavoriteList(t *testing.T) {
	token, _ := auth.GenerateToken("chy", 7777)
	url := "http://127.0.0.1:8080/douyin/favorite/list/?user_id=7777&token=" + token
	method := "GET"
	SendRequest(method, url, nil)
}
