package controller

import (
	"strings"
	"testing"

	"github.com/Shanwu404/TikTokLite/middleware/auth"
)

func TestRelationAction(t *testing.T) {
	auth_token, _ := auth.GenerateToken("Lihua", 2)
	token := "token=" + auth_token
	url := "http://localhost:8080/douyin/relation/action/?to_user_id=7&action_type=1"
	method := "POST"
	SendRequest(method, url, strings.NewReader(token))
}

func TestFollowerList(t *testing.T) {
	token, _ := auth.GenerateToken("Lihua", 2)
	url := "http://localhost:8080/douyin/relation/follower/list/?user_id=7&token=" + token
	method := "GET"
	SendRequest(method, url, strings.NewReader(token))
}
