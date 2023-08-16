package controller

import (
	"strings"
	"testing"

	"github.com/Shanwu404/TikTokLite/utils/auth"
)

func TestRelationAction(t *testing.T) {
	auth_token, _ := auth.GenerateToken("Lqs1", 10)
	token := "token=" + auth_token
	url := "http://127.0.0.1:8080/douyin/relation/action/?to_user_id=600&action_type=2"
	method := "POST"
	SendRequest(method, url, strings.NewReader(token))
}
