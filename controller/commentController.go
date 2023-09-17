package controller

import (
	"net/http"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/gin-gonic/gin"
)

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
		c.JSON(http.StatusBadRequest, facade.CommentActionResponse{
			Response: facade.Response{StatusCode: -1, StatusMsg: "Invalid Request."},
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
		c.JSON(http.StatusBadRequest, facade.CommentListResponse{
			Response: facade.Response{StatusCode: -1, StatusMsg: "Invalid Request."},
		})
		return
	}
	c.JSON(http.StatusOK, cc.commentFacade.CommentList(req))
	return
}
