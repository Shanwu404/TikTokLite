package controller

import (
	"net/http"
	"time"

	"github.com/Shanwu404/TikTokLite/facade"

	"github.com/gin-gonic/gin"
)

const (
	FeedLimit = 3
)

type VideoController struct {
	facade facade.VideoFacade
}

func NewVideoController() *VideoController {
	return &VideoController{
		facade.NewVideoFacade(),
	}
}

// Feed GET /douyin/feed/ 视频流接口
func (vc *VideoController) Feed(c *gin.Context) {
	reqParams, _ := feedParseAndValidateParams(c)
	userId := c.GetInt64("id")
	latestTime := reqParams.LatestTime
	if latestTime == 0 {
		latestTime = time.Now().Unix()
		reqParams.LatestTime = latestTime
	}

	c.JSON(http.StatusOK, *vc.facade.Feed(&reqParams, userId))
}

func (vc *VideoController) PublishAction(c *gin.Context) {
	req, valid := publishActionParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, facade.DouyinPublishActionResponse{
			Response: facade.Response{StatusCode: -1, StatusMsg: "Invalid Request."},
		})
		return
	}
	httpStatus, publishAcyionResponse := vc.facade.PublishAction(&req, c.GetString("username"))
	c.JSON(httpStatus, *publishAcyionResponse)
}

func (vc *VideoController) PublishList(c *gin.Context) {
	reqParams, valid := publishListParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, facade.DouyinPublishListResponse{
			Response: facade.Response{StatusCode: -1, StatusMsg: "Invalid Request."},
		})
		return
	}

	userId := c.GetInt64("id")
	httpStatus, publishListResponse := vc.facade.PublishList(&reqParams, userId)
	c.JSON(httpStatus, publishListResponse)
}
