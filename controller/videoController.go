package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/Shanwu404/TikTokLite/service"

	"github.com/gin-gonic/gin"
)

const (
	FeedLimit = 8
)

type VideoController struct {
	videoService    service.VideoService
	userService     service.UserService
	commentService  service.CommentService
	likeService     service.LikeService
	relationService service.RelationService
}

func NewVideoController() *VideoController {
	return &VideoController{
		service.NewVideoService(),
		service.NewUserService(),
		service.NewCommentService(),
		service.NewLikeService(),
		service.NewRelationService(),
	}
}

// Feed GET /douyin/feed/ 视频流接口
func (vc *VideoController) Feed(c *gin.Context) {
	reqParams, _ := feedParseAndValidateParams(c)
	userId := c.GetInt64("id")
	latestTime := reqParams.LatestTime
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

	//目前客户端可缓存30个视频
	videosWithAuthorID := vc.videoService.GetMultiVideoBefore(latestTime, FeedLimit)
	nextTimeInt := time.Now().Unix()
	if len(videosWithAuthorID) == FeedLimit && FeedLimit > 0 {
		nextTimeInt = videosWithAuthorID[len(videosWithAuthorID)-1].PublishTime.Unix()
	}
	if len(videosWithAuthorID) < FeedLimit {
		nextTimeInt = 0
	}
	videoList := make([]Video, len(videosWithAuthorID))

	// TODO: Goroutine
	for i := range videosWithAuthorID {
		authorInfo := UserInfo{Id: videosWithAuthorID[i].AuthorID}
		vc.completeUserInfo(&authorInfo, userId)
		vc.combineVideoAndAuthor(&videosWithAuthorID[i], &authorInfo, &videoList[i], userId)
	}
	c.JSON(http.StatusOK, douyinFeedResponse{
		Response:  Response{0, "Feeding Succeeded."},
		NextTime:  nextTimeInt,
		VideoList: videoList,
	})
}

func (vc *VideoController) PublishAction(c *gin.Context) {
	req, valid := publishActionParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, douyinPublishActionResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}
	filename := c.GetString("username") + "_" + req.Title + "_" + time.Now().Format("20060102150405")
	// err := c.SaveUploadedFile(req.Data, "videos/"+filename)
	author, err := vc.userService.QueryUserByUsername(c.GetString("username"))
	if err != nil {
		log.Println("[PublishAction]Error when query userID:", err)
		c.JSON(http.StatusInternalServerError, douyinPublishActionResponse{
			Response: Response{2, "Recording Failed."},
		})
		return
	}
	video := service.VideoParams{
		AuthorID:    author.ID,
		PlayURL:     "",
		CoverURL:    "",
		PublishTime: time.Now().Truncate(time.Second),
		Title:       req.Title,
	}
	err = vc.videoService.StoreVideo(req.Data, filename, &video)
	if err != nil {
		log.Println("Uploading failed:" + err.Error())
		c.JSON(http.StatusInternalServerError, douyinPublishActionResponse{
			Response: Response{1, "Uploading Failed."},
		})
		return
	}
	c.JSON(http.StatusOK, douyinPublishActionResponse{
		Response: Response{0, "Uploaded."},
	})
}

func (vc *VideoController) PublishList(c *gin.Context) {
	reqParams, valid := publishListParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, douyinPublishListResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}
	userId := c.GetInt64("id")
	userWorks := vc.videoService.GetVideoListByUserId(reqParams.UserID)
	authorInfo := UserInfo{Id: userWorks[0].AuthorID} // 同一个用户的视频，所以作者信息是一样的
	vc.completeUserInfo(&authorInfo, userId)
	videoList := make([]Video, len(userWorks))
	for i := range userWorks {
		vc.combineVideoAndAuthor(&userWorks[i], &authorInfo, &videoList[i], userId)
	}
	c.JSON(http.StatusOK, douyinPublishListResponse{
		Response:  Response{0, "Get Publish List."},
		VideoList: videoList,
	})
}

// --------------------------------
// 这部分工具函数也要跟随组装数据代码一起放入单独一层

func (vc *VideoController) combineVideoAndAuthor(video *service.VideoParams, author *UserInfo, result *Video, userId int64) {
	flag, _ := vc.likeService.IsLike(video.ID, userId)
	*result = Video{
		ID:            video.ID,
		Author:        *author,
		PlayURL:       video.PlayURL,
		CoverURL:      video.CoverURL,
		FavoriteCount: vc.likeService.CountLikes(video.ID),
		CommentCount:  vc.commentService.CountComments(video.ID),
		IsFavorite:    flag,
		Title:         video.Title,
	}
}

func (vc *VideoController) completeUserInfo(userinfo *UserInfo, userId int64) {
	brief, _ := vc.userService.QueryUserByID(userinfo.Id)
	follows, _ := vc.relationService.CountFollows(brief.ID)
	followers, _ := vc.relationService.CountFollowers(brief.ID)
	isFollow, _ := vc.relationService.IsFollowed(userId, brief.ID)
	favorite_count, _ := vc.likeService.LikeVideoCount(brief.ID)
	*userinfo = UserInfo{
		Id:              brief.ID,
		Username:        brief.Username,
		FollowCount:     follows,   // followService提供
		FollowerCount:   followers, // followService提供
		IsFollow:        isFollow,  // followService提供
		Avatar:          "https://image.zhihuishu.com/zhs/ablecommons/demo/201804/a3b5f5570a2740749d3c372848a18d6f.jpg",
		BackgroundImage: "https://image.zhihuishu.com/zhs/ablecommons/demo/201804/a3b5f5570a2740749d3c372848a18d6f.jpg",
		Signature:       "唯一不变是永远的改变",
		TotalFavorited:  vc.likeService.TotalFavorited(brief.ID), // likeService提供
		WorkCount:       int64(len(vc.videoService.GetVideoListByUserId(userinfo.Id))),
		FavoriteCount:   favorite_count, // likeService提供
	}
}
