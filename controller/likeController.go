package controller

import (
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteListResponse struct {
	Response
	VideoInfo []VideoInfo
}

type LikeController struct {
	videoService    service.VideoService
	userService     service.UserService
	relationService service.RelationService
	likeService     service.LikeService
}

func NewLikeController() *LikeController {
	return &LikeController{
		videoService:    service.NewVideoService(),
		userService:     service.NewUserService(),
		relationService: service.NewRelationService(),
		likeService:     service.NewLikeSerivce(),
	}
}

// FavoriteAction POST /douyin/favorite/action/ 赞操作
func (lc *LikeController) FavoriteAction(c *gin.Context) {
	lsi := service.LikeServiceImpl{}

	actionType := c.Query("action_type")
	userId := c.GetInt64("id")
	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)
	if actionType == "1" {
		if err := lsi.Like(userId, videoId); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "failed"})
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
	} else if actionType == "2" {
		if err := lsi.Unlike(userId, videoId); err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "failed"})
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "success"})
	}
}

// FavoriteList GET /douyin/favorite/list/ 喜欢列表
func (lc *LikeController) FavoriteList(c *gin.Context) {
	vsi := service.VideoServiceImpl{}
	usi := service.UserServiceImpl{}
	rsi := service.RelationServiceImpl{}
	lsi := service.LikeServiceImpl{}

	id := c.Query("user_id")
	userId, _ := strconv.ParseInt(id, 10, 64)
	videos := vsi.GetVideoListByUserId(userId)
	var videoList []VideoInfo
	for _, videoParam := range videos {
		user, _ := usi.QueryUserByID(videoParam.AuthorID)
		followCount, _ := rsi.CountFollows(user.ID)
		followerCount, _ := rsi.CountFollowers(user.ID)
		isFollowed, _ := rsi.IsFollowed(userId, user.ID)
		favoriteCount, _ := lsi.LikeVideoCount(user.ID)
		videoList = append(videoList, VideoInfo{
			User: UserInfo{
				Id:            user.ID,
				Username:      user.Username,
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      isFollowed,
				WorkCount:     int64(len(vsi.GetVideoListByUserId(user.ID))),
				FavoriteCount: favoriteCount,
			},
		})
	}
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoInfo: videoList,
	})
	return
}
