package facade

import (
	"sync"
	"time"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/Shanwu404/TikTokLite/utils"
)

type CommentInfo struct {
	Id         int64                  `json:"id"`
	User       service.UserInfoParams `json:"user"`
	Content    string                 `json:"content"`
	CreateDate string                 `json:"create_date"`
}

type CommentActionRequest struct {
	UserId     int64
	VideoId    int64
	ActionType string
	Content    string
	CommentId  int64
}

type CommentActionResponse struct {
	Response
	CommentInfo CommentInfo `json:"comment,omitempty"`
}

type CommentListRequest struct {
	Video service.VideoParams
}

type CommentListResponse struct {
	Response
	CommentList []CommentInfo `json:"comment_list"`
}

type CommentFacade struct {
	userService     service.UserService
	commentService  service.CommentService
	relationService service.RelationService
	videoService    service.VideoService
	likeService     service.LikeService
}

func NewCommentFacade() *CommentFacade {
	return &CommentFacade{
		userService:     service.NewUserService(),
		commentService:  service.NewCommentService(),
		relationService: service.NewRelationService(),
		videoService:    service.NewVideoService(),
		likeService:     service.NewLikeService(),
	}
}

func (cf *CommentFacade) CommentAction(req CommentActionRequest) CommentActionResponse {
	if req.ActionType == "1" {
		content := utils.Filter.Replace(req.Content, '#')
		// 获取当前视频
		video := cf.videoService.QueryVideoById(req.VideoId)
		t := time.Now()
		comment := service.CommentParams{
			UserId:     req.UserId,
			VideoId:    req.VideoId,
			Content:    content,
			CreateDate: t,
		}
		commentId, code, message := cf.commentService.PostComment(comment)
		if code != 0 {
			return CommentActionResponse{
				Response: Response{StatusCode: code, StatusMsg: message},
			}
		}
		userInfo, _ := cf.userService.QueryUserInfoByID(req.UserId)
		isFollow, _ := cf.relationService.IsFollowed(req.UserId, video.AuthorID)
		userInfo.IsFollow = isFollow
		return CommentActionResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			CommentInfo: CommentInfo{
				Id:         commentId,
				User:       userInfo,
				Content:    content,
				CreateDate: utils.TimeToStr(t),
			},
		}
	} else {
		code, message := cf.commentService.DeleteComment(req.CommentId)
		return CommentActionResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		}
	}
}

func (cf *CommentFacade) CommentList(req CommentListRequest) CommentListResponse {
	comments := cf.commentService.QueryCommentsByVideoId(req.Video.ID)
	commentList := make([]CommentInfo, len(comments))
	var wg sync.WaitGroup
	wg.Add(len(comments))
	for idx, comment := range comments {
		go func(idx int, comment service.CommentParams) {
			defer wg.Done()
			user, err := cf.userService.QueryUserByID(comment.UserId)
			if err != nil {
				return
			}
			userInfo, _ := cf.userService.QueryUserInfoByID(user.ID)
			isFollow, _ := cf.relationService.IsFollowed(user.ID, req.Video.AuthorID)
			userInfo.IsFollow = isFollow
			commentList[idx] = CommentInfo{
				Id:         comment.Id,
				User:       userInfo,
				Content:    comment.Content,
				CreateDate: utils.TimeToStr(comment.CreateDate),
			}
		}(idx, comment)
	}
	wg.Wait()
	return CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: commentList,
	}
}
