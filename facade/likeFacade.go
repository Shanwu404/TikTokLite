package facade

import (
	"github.com/Shanwu404/TikTokLite/service"
)

type LikeActionRequest struct {
	UserId     int64
	VideoId    int64
	ActionType string
}

type LikeActionResponse struct {
	Response
	// VideoList []service.VideoParams
	// VideoList []Video `json:"video_list,omitempty"`
}

type LikeListRequest struct {
	UserId int64
}

type LikeListResponse struct {
	Response
	// VideoList []service.VideoParams
	VideoList []Video `json:"video_list,omitempty"`
}

type LikeFacade struct {
	userService     service.UserService
	commentService  service.CommentService
	relationService service.RelationService
	videoService    service.VideoService
	likeService     service.LikeService
}

func NewLikeFacade() *LikeFacade {
	return &LikeFacade{
		userService:     service.NewUserService(),
		commentService:  service.NewCommentService(),
		relationService: service.NewRelationService(),
		videoService:    service.NewVideoService(),
		likeService:     service.NewLikeService(),
	}
}

func (lf *LikeFacade) FavoriteAction(req LikeActionRequest) LikeActionResponse {

	if req.ActionType == "1" {
		if err := lf.likeService.Like(req.UserId, req.VideoId); err != nil {
			return LikeActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "failed"},
			}
		}
		return LikeActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "like success"},
		}
	} else {
		if err := lf.likeService.Unlike(req.UserId, req.VideoId); err != nil {
			return LikeActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "failed"},
			}
		}
		return LikeActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "unlike success"},
		}
	}

}

// FavoriteList GET /douyin/favorite/list/ 喜欢列表
func (lf *LikeFacade) FavoriteList(req LikeListRequest) LikeListResponse {

	videos := lf.likeService.GetLikeLists(req.UserId)
	var videoList = make([]Video, 0, len(videos))
	for _, videoParam := range videos { //将视频的详细信息格式化

		authorInfo, _ := lf.userService.QueryUserInfoByID(videoParam.AuthorID)
		// authorInfo := UserInfo{Id: videoParam.AuthorID}
		// vc.completeUserInfo(&authorInfo, userId)
		authorInfo.FollowCount, _ = lf.relationService.CountFollows(authorInfo.Id)
		authorInfo.FollowerCount, _ = lf.relationService.CountFollowers(authorInfo.Id)
		authorInfo.IsFollow, _ = lf.relationService.IsFollowed(req.UserId, authorInfo.Id)
		authorInfo.WorkCount = int64(len(lf.videoService.GetVideoListByUserId(authorInfo.Id)))
		authorInfo.FavoriteCount = lf.likeService.CountLikes(authorInfo.Id)
		video := Video{
			ID:            videoParam.ID,
			Author:        authorInfo,
			PlayURL:       videoParam.PlayURL,
			CoverURL:      videoParam.CoverURL,
			FavoriteCount: lf.likeService.CountLikes(videoParam.ID),
			CommentCount:  lf.commentService.CountComments(videoParam.ID),
			IsFavorite:    lf.likeService.IsLike(videoParam.ID, req.UserId),
			Title:         videoParam.Title,
		}
		videoList = append(videoList, video)
	}
	return LikeListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "success"},
		VideoList: videoList,
	}
}
