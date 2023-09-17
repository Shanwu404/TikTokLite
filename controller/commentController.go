package controller

import (
	"net/http"

	"github.com/Shanwu404/TikTokLite/facade"
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
	commentFacade facade.CommentFacade
}

func NewCommentController() *CommentController {
	return &CommentController{
		commentFacade: *facade.NewCommentFacade(),
	}
}

// CommentAction POST /douyin/comment/action/ 评论操作
func (cc *CommentController) CommentAction(c *gin.Context) {
	req, valid := CommentActionParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, CommentActionResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}
	c.JSON(http.StatusOK, cc.commentFacade.CommentAction(req))
	return
}

// CommentList GET /douyin/comment/list/ 评论列表
func (cc *CommentController) CommentList(c *gin.Context) {
	req, valid := CommentListParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, CommentListResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}
	c.JSON(http.StatusOK, cc.commentFacade.CommentList(req))
	return
}
