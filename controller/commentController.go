package controller

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []CommentInfo
}

type CommentActionResponse struct {
	Response
	CommentInfo CommentInfo
}

// CommentAction POST /douyin/comment/action/ 评论操作
func CommentAction(c *gin.Context) {
	csi := service.CommentServiceImpl{}
	usi := service.UserServiceImpl{}

	actionType := c.Query("action_type")
	if actionType == "1" {
		content := c.Query("comment_text")
		// 获取当前用户
		userId := c.GetInt64("id")
		// 获取当前视频
		id := c.Query("video_id")
		videoId, _ := strconv.ParseInt(id, 10, 64)
		t := time.Now()
		comment := dao.Comment{
			UserId:     userId,
			VideoId:    videoId,
			Content:    content,
			CreateDate: t,
		}
		commentId, code, message := csi.PostComment(comment)
		if code != 0 {
			c.JSON(http.StatusOK, Response{
				StatusCode: code,
				StatusMsg:  message,
			})
			return
		}
		userInfo, _ := usi.QueryUserRespByID(userId)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			CommentInfo: CommentInfo{
				Id:         commentId,
				User:       userInfo,
				Content:    content,
				CreateDate: time.Now(),
			},
		})
		return
	} else {
		cId := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(cId, 10, 64)
		code, message := csi.DeleteComment(commentId)
		c.JSON(http.StatusOK, Response{StatusCode: code, StatusMsg: message})
		return
	}
}

// CommentList GET /douyin/comment/list/ 评论列表
func CommentList(c *gin.Context) {
	usi := service.UserServiceImpl{}
	csi := service.CommentServiceImpl{}

	id := c.Query("video_id")
	videoId, _ := strconv.ParseInt(id, 10, 64)
	comments := csi.QueryCommentsByVideoId(videoId)
	var commonList []CommentInfo
	for _, comment := range comments {
		user, err := usi.QueryUserRespByID(comment.UserId)
		if err != nil {
			continue
		}
		commonList = append(commonList, CommentInfo{
			Id:         comment.Id,
			User:       user,
			Content:    comment.Content,
			CreateDate: time.Now(),
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: commonList,
	})
	return
}
