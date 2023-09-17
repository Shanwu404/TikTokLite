package facade

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/service"
)

type DouyinFeedRequest struct {
	LatestTime int64  `form:"latest_time"`
	Token      string `form:"token"`
}

type DouyinFeedResponse struct {
	Response
	NextTime  int64   `json:"next_time,omitempty"`
	VideoList []Video `json:"video_list"`
}

type DouyinPublishActionRequest struct {
	Token string                `json:"token"`
	Data  *multipart.FileHeader `json:"data"`
	Title string                `json:"tilte"`
}

type DouyinPublishActionResponse struct {
	Response
}

type DouyinPublishListRequest struct {
	UserID int64  `form:"user_id"`
	Token  string `form:"token"`
}

type DouyinPublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

type Video struct {
	ID            int64                  `json:"id"`
	Author        service.UserInfoParams `json:"author"`
	PlayURL       string                 `json:"play_url"`
	CoverURL      string                 `json:"cover_url"`
	FavoriteCount int64                  `json:"favorite_count"`
	CommentCount  int64                  `json:"comment_count"`
	IsFavorite    bool                   `json:"is_favorite"`
	Title         string                 `json:"title"`
}

const (
	FeedLimit = 3
)

type VideoFacade struct {
	videoService    service.VideoService
	userService     service.UserService
	commentService  service.CommentService
	likeService     service.LikeService
	relationService service.RelationService
}

func NewVideoFacade() VideoFacade {
	return VideoFacade{
		service.NewVideoService(),
		service.NewUserService(),
		service.NewCommentService(),
		service.NewLikeService(),
		service.NewRelationService(),
	}
}

// Feed GET /douyin/feed/ 视频流接口
func (vf *VideoFacade) Feed(req *DouyinFeedRequest, userId int64) *DouyinFeedResponse {
	//目前客户端可缓存30个视频
	videosWithAuthorID := vf.videoService.GetMultiVideoBefore(req.LatestTime, FeedLimit)
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
		authorInfo := service.UserInfoParams{Id: videosWithAuthorID[i].AuthorID}
		vf.completeUserInfo(&authorInfo, userId)
		vf.combineVideoAndAuthor(&videosWithAuthorID[i], &authorInfo, &videoList[i], userId)
	}
	return &DouyinFeedResponse{
		Response:  Response{0, "Feeding Succeeded."},
		NextTime:  nextTimeInt,
		VideoList: videoList,
	}
}

func (vf *VideoFacade) PublishAction(req *DouyinPublishActionRequest, username string) (int, *DouyinPublishActionResponse) {
	filename := username + "_" + req.Title + "_" + time.Now().Format("20060102150405")
	author, err := vf.userService.QueryUserByUsername(username)
	if err != nil {
		logger.Errorln("[PublishAction]Error when query userID:", err)
		return http.StatusInternalServerError, &DouyinPublishActionResponse{
			Response: Response{2, "Recording Failed."},
		}
	}
	video := service.VideoParams{
		AuthorID:    author.ID,
		PlayURL:     "",
		CoverURL:    "",
		PublishTime: time.Now().Truncate(time.Second),
		Title:       req.Title,
	}
	err = vf.videoService.StoreVideo(req.Data, filename, &video)
	if err != nil {
		logger.Errorln("[PublishAction]Uploading failed:", err)
		return http.StatusInternalServerError, &DouyinPublishActionResponse{
			Response: Response{1, "Uploading Failed."},
		}
	}
	return http.StatusOK, &DouyinPublishActionResponse{
		Response: Response{0, "Uploaded."},
	}
}

func (vf *VideoFacade) PublishList(req *DouyinPublishListRequest, userId int64) (int, *DouyinPublishListResponse) {
	userWorks := vf.videoService.GetVideoListByUserId(req.UserID)
	if len(userWorks) == 0 {
		return http.StatusOK, &DouyinPublishListResponse{
			Response: Response{0, "Get Publish List."},
		}
	}
	authorInfo := service.UserInfoParams{Id: userWorks[0].AuthorID} // 同一个用户的视频，所以作者信息是一样的
	vf.completeUserInfo(&authorInfo, userId)
	videoList := make([]Video, len(userWorks))
	for i := range userWorks {
		vf.combineVideoAndAuthor(&userWorks[i], &authorInfo, &videoList[i], userId)
	}
	return http.StatusOK, &DouyinPublishListResponse{
		Response:  Response{0, "Get Publish List."},
		VideoList: videoList,
	}
}

// --------------------------------

func (vf *VideoFacade) combineVideoAndAuthor(video *service.VideoParams, author *service.UserInfoParams, result *Video, userId int64) {
	flag := vf.likeService.IsLike(video.ID, userId)
	*result = Video{
		ID:            video.ID,
		Author:        *author,
		PlayURL:       video.PlayURL,
		CoverURL:      video.CoverURL,
		FavoriteCount: vf.likeService.CountLikes(video.ID),
		CommentCount:  vf.commentService.CountComments(video.ID),
		IsFavorite:    flag,
		Title:         video.Title,
	}
}

func (vf *VideoFacade) completeUserInfo(userinfo *service.UserInfoParams, userId int64) {
	brief, _ := vf.userService.QueryUserByID(userinfo.Id)
	follows, _ := vf.relationService.CountFollows(brief.ID)
	followers, _ := vf.relationService.CountFollowers(brief.ID)
	isFollow, _ := vf.relationService.IsFollowed(userId, brief.ID)
	favorite_count, _ := vf.likeService.LikeVideoCount(brief.ID)
	*userinfo = service.UserInfoParams{
		Id:              brief.ID,
		Username:        brief.Username,
		FollowCount:     follows,   // followService提供
		FollowerCount:   followers, // followService提供
		IsFollow:        isFollow,  // followService提供
		Avatar:          "https://mary-aliyun-img.oss-cn-beijing.aliyuncs.com/typora/202308171029672.jpg",
		BackgroundImage: "https://mary-aliyun-img.oss-cn-beijing.aliyuncs.com/typora/202308171029672.jpg",
		Signature:       "这个人很懒，什么都没有留下",
		TotalFavorited:  vf.likeService.TotalFavorited(brief.ID), // likeService提供
		WorkCount:       int64(len(vf.videoService.GetVideoListByUserId(userinfo.Id))),
		FavoriteCount:   favorite_count, // likeService提供
	}
}
