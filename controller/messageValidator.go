package controller

import (
	"strconv"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

func MessageActionParseAndValidateParams(c *gin.Context) (facade.MessageActionRequest, bool) {
	req := facade.MessageActionRequest{}
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

func MessageListParseAndValidateParams(c *gin.Context) (facade.MessageListRequest, bool) {
	req := facade.MessageListRequest{}
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

	preMsgTimeStr := c.Query("pre_msg_time")
	preMsgTime, _ := strconv.ParseInt(preMsgTimeStr, 10, 64)
	req.PreMsgTime = preMsgTime
	return req, true
}
