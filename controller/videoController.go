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
	videoService service.VideoService
	userService  service.UserService
}

func NewVideoController() *VideoController {
	return &VideoController{
		service.NewVideoService(),
		service.NewUserService(),
	}
}

// Feed GET /douyin/feed/ 视频流接口
func (vc *VideoController) Feed(c *gin.Context) {
	reqParams, _ := feedParseAndValidateParams(c)
	latestTime := reqParams.LatestTime
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

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
	err := c.SaveUploadedFile(req.Data, "videos/"+filename)
	// err := vc.videoService.StoreVideo(req.Data, c.GetString("username"), filename)
	if err != nil {
		log.Println("Uploading failed:" + err.Error())
		c.JSON(http.StatusOK, douyinPublishActionResponse{
			Response: Response{1, "Uploading Failed."},
		})
		return
	}
	author, _ := vc.userService.QueryUserByUsername(c.GetString("username"))
	video := service.VideoParams{
		AuthorID:    author.ID,
		PlayURL:     "http://47.94.162.202:8080/douyin/tiktok/" + filename,
		CoverURL:    "",
		PublishTime: time.Now().Truncate(time.Second),
		Title:       req.Title,
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

func (vc *VideoController) PublishList(c *gin.Context) {
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
	videoList := make([]Video, len(userWorks))
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

func (vc *VideoController) completeUserInfo(userinfo *UserInfo) {
	brief, _ := vc.userService.QueryUserByID(userinfo.Id)
	*userinfo = UserInfo{
		Id:              brief.ID,
		Username:        brief.Username,
		FollowCount:     110,   // followService提供
		FollowerCount:   12000, // followService提供
		IsFollow:        false, // followService提供
		Avatar:          "https://image.zhihuishu.com/zhs/ablecommons/demo/201804/a3b5f5570a2740749d3c372848a18d6f.jpg",
		BackgroundImage: "https://image.zhihuishu.com/zhs/ablecommons/demo/201804/a3b5f5570a2740749d3c372848a18d6f.jpg",
		Signature:       "唯一不变是永远的改变",
		TotalFavorited:  210, // likeService提供
		WorkCount:       int64(len(vc.videoService.GetVideoListByUserId(userinfo.Id))),
		FavoriteCount:   123456, // likeService提供
	}
}
