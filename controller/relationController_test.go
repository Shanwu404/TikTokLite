package controller

import (
	"strings"
	"testing"

	"github.com/Shanwu404/TikTokLite/middleware/auth"
)

func TestRelationAction(t *testing.T) {
	token, _ := auth.GenerateToken("lux", 7)
	url := "http://localhost:8080/douyin/relation/action/?to_user_id=6&action_type=1&token=" + token
	method := "POST"
	SendRequest(method, url, nil)
}

func TestFollowsList(t *testing.T) {
	token, _ := auth.GenerateToken("lux", 7)
	url := "http://localhost:8080/douyin/relation/follow/list/?user_id=7&token=" + token
	method := "GET"
	SendRequest(method, url, strings.NewReader(token))
}

func TestFollowerList(t *testing.T) {
	token, _ := auth.GenerateToken("lux", 7)
	url := "http://localhost:8080/douyin/relation/follower/list/?user_id=7&token=" + token
	method := "GET"
	SendRequest(method, url, strings.NewReader(token))
}

func TestFriendList(t *testing.T) {
	token, _ := auth.GenerateToken("lux", 7)
	url := "http://localhost:8080/douyin/relation/friend/list/?user_id=7&token=" + token
	method := "GET"
	SendRequest(method, url, strings.NewReader(token))
}
