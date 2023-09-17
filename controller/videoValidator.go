// 数据校验
package controller

import (
	"time"
	"unicode/utf8"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/gin-gonic/gin"
)

func feedParseAndValidateParams(c *gin.Context) (facade.DouyinFeedRequest, bool) {
	req := facade.DouyinFeedRequest{}
	c.ShouldBindQuery(&req)
	current := time.Now().Unix()
	if req.LatestTime > current {
		req.LatestTime = current
	}
	return req, true
}

func publishActionParseAndValidateParams(c *gin.Context) (facade.DouyinPublishActionRequest, bool) {
	req := facade.DouyinPublishActionRequest{
		Title: c.PostForm("title"),
	}
	data, err := c.FormFile("data")
	if err != nil {
		return req, false
	}
	req.Data = data
	if req.Data == nil {
		return req, false
	}
	if utf8.RuneCountInString(req.Title) > 255 {
		return req, false
	}
	return req, true
}

func publishListParseAndValidateParams(c *gin.Context) (facade.DouyinPublishListRequest, bool) {
	req := facade.DouyinPublishListRequest{}
	c.ShouldBindQuery(&req)
	if req.UserID <= 0 {
		return req, false
	}
	return req, true
}
