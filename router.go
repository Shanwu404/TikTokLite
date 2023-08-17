package main

import (
	"github.com/Shanwu404/TikTokLite/controller"
	"github.com/Shanwu404/TikTokLite/middleware/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	apiRouter := r.Group("/douyin")
	apiRouter.Static("tiktok", "./videos") //测试用。配置静态资源路径

	userController := controller.NewUserController()
	relationController := controller.NewRelationController()
	videoController := controller.NewVideoController()
	commentController := controller.NewCommentController()
	messageController := controller.NewMessageController()
	likeController := controller.NewLikeController()

	// basic apis
	apiRouter.GET("/feed/", videoController.Feed)
	apiRouter.POST("/publish/action/", auth.Auth, videoController.PublishAction)
	apiRouter.POST("/publish/list/", auth.Auth, videoController.PublishList)

	apiRouter.POST("/user/register/", userController.Register)
	apiRouter.POST("/user/login/", userController.Login)
	apiRouter.GET("/user/", userController.GetUserInfo)

	apiRouter.POST("/comment/action/", auth.Auth, commentController.CommentAction)
	apiRouter.GET("/relation/follow/list/", auth.Auth, relationController.FollowsList)
	apiRouter.GET("/relation/follower/list/", auth.Auth, relationController.FollowersList)
	apiRouter.GET("/comment/list/", auth.Auth, commentController.CommentList)

	apiRouter.POST("/message/action/", auth.Auth, messageController.MessageAction)
	apiRouter.GET("/message/chat/", auth.Auth, messageController.MessageList)

	apiRouter.POST("/favorite/action/", auth.Auth, likeController.FavoriteAction)
	apiRouter.GET("/favorite/list/", auth.Auth, likeController.FavoriteList)
	return r

}
