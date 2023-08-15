// 数据校验
package controller

import (
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type douyinFeedRequest struct {
	LatestTime int64  `form:"latest_time"`
	Token      string `form:"token"`
}

type douyinFeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type douyinPublishActionRequest struct {
	Token string `json:"token"`
	Data  []byte `json:"data"`
	Title string `json:"tilte"`
}

type douyinPublishActionResponse struct {
	Response
}

type douyinPublishListRequest struct {
	UserID int64  `form:"user_id"`
	Token  string `form:"token"`
}

type douyinPublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

type Video struct {
	ID            int64    `json:"id"`
	Author        UserInfo `json:"author"`
	PlayURL       string   `json:"play_url"`
	CoverURL      string   `json:"cover_url"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
}

func feedParseAndValidateParams(c *gin.Context) (douyinFeedRequest, bool) {
	req := douyinFeedRequest{}
	c.ShouldBindQuery(&req)
	current := time.Now().Unix()
	if req.LatestTime > current {
		req.LatestTime = current
	}
	return req, true
}

func publishActionParseAndValidateParams(c *gin.Context) (douyinPublishActionRequest, bool) {
	req := douyinPublishActionRequest{}
	c.ShouldBindJSON(&req)
	if len(req.Data) == 0 {
		return req, false
	}
	if utf8.RuneCountInString(req.Title) > 255 {
		return req, false
	}
	return req, true
}

func publishListParseAndValidateParams(c *gin.Context) (douyinPublishListRequest, bool) {
	req := douyinPublishListRequest{}
	c.ShouldBindQuery(&req)
	if req.UserID <= 0 {
		return req, false
	}
	return req, true
}

func timeStampStr2Time(timeStampStr string) time.Time {
	latestTimeStamp, _ := strconv.Atoi(timeStampStr)
	latestTime := time.Unix(int64(latestTimeStamp), 0)
	return latestTime
}
