package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/Shanwu404/TikTokLite/utils"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []CommentInfo `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	CommentInfo CommentInfo `json:"comment,omitempty"`
}

type CommentController struct {
	userService     service.UserService
	commentService  service.CommentService
	relationService service.RelationService
	videoService    service.VideoService
	likeService     service.LikeService
}

func NewCommentController() *CommentController {
	return &CommentController{
		userService:     service.NewUserService(),
		commentService:  service.NewCommentService(),
		relationService: service.NewRelationService(),
		videoService:    service.NewVideoService(),
		likeService:     service.NewLikeService(),
	}
}

// CommentAction POST /douyin/comment/action/ 评论操作
func (cc *CommentController) CommentAction(c *gin.Context) {
	actionType := c.Query("action_type")
	if actionType == "1" {
		content := c.Query("comment_text")
		// 获取当前用户
		userId := c.GetInt64("id")
		// 获取当前视频
		id := c.Query("video_id")
		videoId, _ := strconv.ParseInt(id, 10, 64)
		video := cc.videoService.QueryVideoById(videoId)
		t := time.Now()
		comment := service.CommentParams{
			UserId:     userId,
			VideoId:    videoId,
			Content:    content,
			CreateDate: t,
		}
		commentId, code, message := cc.commentService.PostComment(comment)
		if code != 0 {
			c.JSON(http.StatusOK, Response{
				StatusCode: code,
				StatusMsg:  message,
			})
			return
		}
		user, _ := cc.userService.QueryUserByID(userId)
		followCount, _ := cc.relationService.CountFollows(user.ID)
		followerCount, _ := cc.relationService.CountFollowers(user.ID)
		isFollowed, _ := cc.relationService.IsFollowed(user.ID, video.AuthorID)
		favoriteCount, _ := cc.likeService.LikeVideoCount(user.ID)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			CommentInfo: CommentInfo{
				Id: commentId,
				User: UserInfo{
					Id:            user.ID,
					Username:      user.Username,
					FollowCount:   followCount,
					FollowerCount: followerCount,
					IsFollow:      isFollowed,
					WorkCount:     int64(len(cc.videoService.GetVideoListByUserId(user.ID))),
					FavoriteCount: favoriteCount,
				},
				Content:    content,
				CreateDate: utils.TimeToStr(t),
			},
		})
		return
	} else {
		cId := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(cId, 10, 64)
		code, message := cc.commentService.DeleteComment(commentId)
		c.JSON(http.StatusOK, Response{StatusCode: code, StatusMsg: message})
		return
	}
}

// CommentList GET /douyin/comment/list/ 评论列表
func (cc *CommentController) CommentList(c *gin.Context) {
	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)
	video := cc.videoService.QueryVideoById(videoId)
	comments := cc.commentService.QueryCommentsByVideoId(videoId)
	var commonList []CommentInfo
	for _, comment := range comments {
		user, err := cc.userService.QueryUserByID(comment.UserId)
		if err != nil {
			continue
		}
		followCount, _ := cc.relationService.CountFollows(user.ID)
		followerCount, _ := cc.relationService.CountFollowers(user.ID)
		isFollowed, _ := cc.relationService.IsFollowed(user.ID, video.AuthorID)
		favoriteCount, _ := cc.likeService.LikeVideoCount(user.ID)
		commonList = append(commonList, CommentInfo{
			Id: comment.Id,
			User: UserInfo{
				Id:            user.ID,
				Username:      user.Username,
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      isFollowed,
				WorkCount:     int64(len(cc.videoService.GetVideoListByUserId(user.ID))),
				FavoriteCount: favoriteCount,
			},
			Content:    comment.Content,
			CreateDate: utils.TimeToStr(comment.CreateDate),
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: commonList,
	})
	return
}
