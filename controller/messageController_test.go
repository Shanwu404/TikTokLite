package controller

import (
	"github.com/Shanwu404/TikTokLite/middleware/auth"
	"strings"
	"testing"
)

func TestMessageAction(t *testing.T) {
	tok, _ := auth.GenerateToken("qly", 1000)
	token := "token=" + tok
	url := "http://127.0.0.1:8080/douyin/message/action/?to_user_id=1001&action_type=1&content=test"
	method := "POST"
	SendRequest(method, url, strings.NewReader(token))
}

func TestMessageList(t *testing.T) {
	tok, _ := auth.GenerateToken("qly", 1000)
	url := "http://127.0.0.1:8080/douyin/message/chat/?to_user_id=1001&token=" + tok
	method := "GET"
	SendRequest(method, url, nil)
}
