package controller

import (
	"strconv"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/gin-gonic/gin"
)

func RelationActionParseAndValidateParams(c *gin.Context) (facade.RelationActionRequest, bool) {
	req := facade.RelationActionRequest{}

	var err error

	// 1. 判断to_user_id解析是否有误
	req.ToUserId, err = strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		logger.Infoln("Invalid toUserId:", req.ToUserId)
		return req, false
	}

	// 2. 判断actionType解析是否有误
	req.ActionType, err = strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || req.ActionType < 1 || req.ActionType > 2 {
		logger.Infoln("Invalid actionType:", req.ActionType)
		return req, false
	}

	// 3. 取出用户id
	req.UserId = c.GetInt64("id")

	return req, true
}
