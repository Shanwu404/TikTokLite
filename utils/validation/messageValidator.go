package validation

import (
	"strconv"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type MessageActionRequest struct {
	UserId   int64
	ToUserId int64
	Content  string
}

type MessageListRequest struct {
	UserId   int64
	ToUserId int64
}

func MessageActionParseAndValidateParams(c *gin.Context) (MessageActionRequest, bool) {
	req := MessageActionRequest{}
	usi := service.NewUserService()

	// 判读操作类别
	actionType := c.Query("action_type")
	if actionType != "1" {
		return req, false
	}

	// 判断 id
	userId := c.GetInt64("id")
	id := c.Query("to_user_id")
	toUserId, _ := strconv.ParseInt(id, 10, 64)
	if !usi.IsUserIdExist(userId) {
		return req, false
	}
	req.UserId = userId
	if !usi.IsUserIdExist(toUserId) {
		return req, false
	}
	req.ToUserId = toUserId

	// 判断 content
	content := c.Query("content")
	if len(content) > 500 {
		return req, false
	}
	req.Content = content
	return req, true
}

func MessageListParseAndValidateParams(c *gin.Context) (MessageListRequest, bool) {
	req := MessageListRequest{}
	usi := service.NewUserService()
	userId := c.GetInt64("id")
	id := c.Query("to_user_id")
	toUserId, _ := strconv.ParseInt(id, 10, 64)
	if !usi.IsUserIdExist(userId) {
		return req, false
	}
	req.UserId = userId
	if !usi.IsUserIdExist(toUserId) {
		return req, false
	}
	req.ToUserId = toUserId
	return req, true
}
