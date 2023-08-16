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
	ListLimit = 8
)

type videoController struct {
	videoService service.VideoService
	userService  service.UserService
}

func NewVideoController() *videoController {
	return &videoController{
		service.NewVideoService(),
		service.NewUserService(),
	}
}

// Feed GET /douyin/feed/ 视频流接口
func (vc *videoController) Feed(c *gin.Context) {
	reqParams, _ := feedParseAndValidateParams(c)
	latestTime := reqParams.LatestTime

	//目前客户端可缓存30个视频
	videosWithAuthorID := vc.videoService.GetMultiVideoBefore(latestTime, FeedLimit)
	nextTimeInt := time.Now().Unix()
	if len(videosWithAuthorID) > 0 {
		nextTimeInt = videosWithAuthorID[len(videosWithAuthorID)-1].PublishTime.Unix()
	}
	videoList := make([]Video, len(videosWithAuthorID))

	// TODO: Goroutine
	for i := range videosWithAuthorID {
		authorInfo := UserInfo{Id: videosWithAuthorID[i].AuthorID}
		vc.completeUserInfo(&authorInfo)
		combineVideoAndAuthor(&videosWithAuthorID[i], &authorInfo, &videoList[i])
	}
	log.Println(videoList)
	c.JSON(http.StatusOK, douyinFeedResponse{
		Response:  Response{0, "Feeding Succeeded."},
		VideoList: videoList,
		NextTime:  nextTimeInt,
	})
}

func (vc *videoController) PublishAction(c *gin.Context) {
	reqJSON, valid := publishActionParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, douyinPublishActionResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}
	filename := reqJSON.Title
	err := vc.videoService.StoreVideo(reqJSON.Data, c.GetString("username"), filename)
	if err != nil {
		log.Println("Uploading failed:" + err.Error())
		c.JSON(http.StatusOK, douyinPublishActionResponse{
			Response: Response{1, "Uploading Failed."},
		})
		return
	}
	author, _ := vc.userService.QueryUserByUsername(c.GetString("username"))
	prefix := ""
	video := service.VideoParams{
		AuthorID:    author.ID,
		PlayURL:     prefix + filename,
		CoverURL:    "",
		PublishTime: time.Now().Truncate(time.Second),
		Title:       reqJSON.Title,
	}
	err = vc.videoService.InsertVideosTable(&video)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, douyinPublishActionResponse{
			Response: Response{2, "Recording Failed."},
		})
		return
	}
}

func (vc *videoController) PublishList(c *gin.Context) {
	reqParams, valid := publishListParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, douyinPublishListResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}

	userWorks := vc.videoService.GetVideoListByUserId(reqParams.UserID)
	authorInfo := UserInfo{Id: userWorks[0].AuthorID} // 同一个用户的视频，所以作者信息是一样的
	vc.completeUserInfo(&authorInfo)
	videoList := make([]Video, 0, len(userWorks))
	for i := range userWorks {
		combineVideoAndAuthor(&userWorks[i], &authorInfo, &videoList[i])
	}
	c.JSON(http.StatusOK, douyinPublishListResponse{
		Response:  Response{0, "Get Publish List."},
		VideoList: videoList,
	})

}

// --------------------------------
// 这部分工具函数也要跟随组装数据代码一起放入单独一层

func combineVideoAndAuthor(video *service.VideoParams, author *UserInfo, result *Video) {
	*result = Video{
		ID:            video.ID,
		Author:        *author,
		PlayURL:       video.PlayURL,
		CoverURL:      video.CoverURL,
		FavoriteCount: 1000,
		CommentCount:  1000,
		IsFavorite:    true,
		Title:         video.Title,
	}
}

func (vc *videoController) completeUserInfo(userinfo *UserInfo) {
	brief, _ := vc.userService.QueryUserByID(userinfo.Id)
	*userinfo = UserInfo{
		Id:              brief.ID,
		Username:        brief.Username,
		FollowCount:     110,   // followService提供
		FollowerCount:   12000, // followService提供
		IsFollow:        false, // followService提供
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "唯一不变是永远的改变",
		TotalFavorited:  9876543210, // likeService提供
		WorkCount:       int64(len(vc.videoService.GetVideoListByUserId(userinfo.Id))),
		FavoriteCount:   123456, // likeService提供
	}
}
