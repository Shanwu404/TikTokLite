package controller

import (
	"strconv"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type CommentActionRequest struct {
	UserId     int64
	VideoId    int64
	ActionType string
	Content    string
	CommentId  int64
}

type CommentListRequest struct {
	Video service.VideoParams
}

func CommentActionParseAndValidateParams(c *gin.Context) (CommentActionRequest, bool) {
	req := CommentActionRequest{}
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
	if actionType == "1" {
		content := c.Query("comment_text")
		if len(content) > 500 {
			return req, false
		}
		req.Content = content
		return req, true
	} else {
		cId := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(cId, 10, 64)
		req.CommentId = commentId
		return req, true
	}
}

func CommentListParseAndValidateParams(c *gin.Context) (CommentListRequest, bool) {
	req := CommentListRequest{}
	vsi := service.NewVideoService()

	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)
	video := vsi.QueryVideoById(videoId)
	if !vsi.Exist(videoId) {
		return req, false
	}
	req.Video = video
	return req, true
}
