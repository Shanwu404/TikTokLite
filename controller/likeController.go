package controller

import (
	"net/http"
	"strconv"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

type LikeController struct {
	videoService    service.VideoService
	userService     service.UserService
	relationService service.RelationService
	likeService     service.LikeService
	commentService  service.CommentService
}

func NewLikeController() *LikeController {
	return &LikeController{
		videoService:    service.NewVideoService(),
		userService:     service.NewUserService(),
		relationService: service.NewRelationService(),
		likeService:     service.NewLikeService(),
		commentService:  service.NewCommentService(),
	}
}

// FavoriteAction POST /douyin/favorite/action/ 赞操作
func (lc *LikeController) FavoriteAction(c *gin.Context) {
	lsi := service.LikeServiceImpl{}

	actionType := c.Query("action_type")
	userId := c.GetInt64("id")
	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)

	// if !redis.CuckooFilterVideoId.Contain([]byte(strconv.FormatInt(videoId, 10))) {
	// 	c.JSON(http.StatusOK, likeResponse{StatusCode: 1, StatusMsg: "视频不存在！"})
	// 	return
	// }

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
	var videoList []Video

	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, FavoriteListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Invalid value."},
		})
		return
	}
	videos := lc.likeService.GetLikeLists(userId)
	// var videoList = make([]Video, 0, len(videos))

	for _, videoParam := range videos { //将视频的详细信息格式化

		authorInfo, _ := lc.userService.QueryUserInfoByID(videoParam.AuthorID)
		// authorInfo := UserInfo{Id: videoParam.AuthorID}
		// vc.completeUserInfo(&authorInfo, userId)
		authorInfo.FollowCount, _ = lc.relationService.CountFollows(authorInfo.Id)
		authorInfo.FollowerCount, _ = lc.relationService.CountFollowers(authorInfo.Id)
		authorInfo.IsFollow, _ = lc.relationService.IsFollowed(userId, authorInfo.Id)
		authorInfo.WorkCount = int64(len(lc.videoService.GetVideoListByUserId(authorInfo.Id)))
		authorInfo.FavoriteCount = lc.likeService.CountLikes(authorInfo.Id)
		video := Video{
			ID:            videoParam.ID,
			Author:        UserInfo(authorInfo),
			PlayURL:       videoParam.PlayURL,
			CoverURL:      videoParam.CoverURL,
			FavoriteCount: lc.likeService.CountLikes(videoParam.ID),
			CommentCount:  lc.commentService.CountComments(videoParam.ID),
			IsFavorite:    lc.likeService.IsLike(videoParam.ID, userId),
			Title:         videoParam.Title,
		}
		videoList = append(videoList, video)
	}
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videoList,
	})
}
