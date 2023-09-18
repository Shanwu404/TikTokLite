package controller

import (
	"strconv"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

func LikeActionParseAndValidateParams(c *gin.Context) (facade.LikeActionRequest, bool) {
	req := facade.LikeActionRequest{}
	usi := service.NewUserService()
	vsi := service.NewVideoService()

	// 判断 userId 是否存在
	userId := c.GetInt64("id")
	flag := usi.IsUserIdExist(userId)
	if !flag {
		return req, false
	}
	req.UserId = userId

	// 判断 videoId 是否存在
	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)
	flag2 := vsi.Exist(videoId)
	if !flag2 {
		return req, false
	}
	req.VideoId = videoId

	// 判读操作类别
	actionType := c.Query("action_type")
	if actionType != "1" && actionType != "2" {
		return req, false
	}
	req.ActionType = actionType
	return req, true

}

func LikeListParseAndValidateParams(c *gin.Context) (facade.LikeListRequest, bool) {
	req := facade.LikeListRequest{}
	usi := service.NewUserService()

	// 判断 userId 是否存在
	id := c.Query("user_id")
	userId, _ := strconv.ParseInt(id, 10, 64)
	// userId := c.GetInt64("user_id")
	flag := usi.IsUserIdExist(userId)
	if !flag {
		return req, false
	}
	req.UserId = userId
	return req, true
}
