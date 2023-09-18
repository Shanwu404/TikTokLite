package controller

import (
	"net/http"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/gin-gonic/gin"
)

type LikeController struct {
	likeFacade facade.LikeFacade
}

func NewLikeController() *LikeController {
	return &LikeController{
		likeFacade: *facade.NewLikeFacade(),
	}
}

// FavoriteAction POST /douyin/favorite/action/ 赞操作
func (lc *LikeController) FavoriteAction(c *gin.Context) {
	req, valid := LikeActionParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, facade.LikeActionResponse{
			Response: facade.Response{StatusCode: -1, StatusMsg: "Invalid Request."},
		})
		return
	}
	c.JSON(http.StatusOK, lc.likeFacade.FavoriteAction(req))
}

// FavoriteList GET /douyin/favorite/list/ 喜欢列表
func (lc *LikeController) FavoriteList(c *gin.Context) {
	req, valid := LikeListParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, facade.LikeListResponse{
			Response: facade.Response{StatusCode: -1, StatusMsg: "Invalid Request."},
		})
		return
	}
	c.JSON(http.StatusOK, lc.likeFacade.FavoriteList(req))

}
