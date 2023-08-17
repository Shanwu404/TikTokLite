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
	videoController := controller.NewVideoController()
	relationController := controller.NewRelationController()
	// basic apis
	apiRouter.GET("/feed/", videoController.Feed)
	apiRouter.POST("/publish/action/", auth.Auth, videoController.PublishAction)
	apiRouter.GET("/publish/list/", auth.Auth, videoController.PublishList)

	apiRouter.POST("/user/register/", userController.Register)
	apiRouter.POST("/user/login/", userController.Login)
	apiRouter.GET("/user/", auth.Auth, userController.GetUserInfo)

	apiRouter.POST("/relation/action/", auth.Auth, relationController.RelationAction)
	apiRouter.GET("/relation/follower/list/", auth.Auth, relationController.FollowersList)

	apiRouter.POST("/comment/action/", auth.Auth, controller.CommentAction)
	apiRouter.GET("/comment/list/", auth.Auth, controller.CommentList)

	apiRouter.POST("/message/action/", auth.Auth, controller.MessageAction)
	apiRouter.GET("/message/chat/", auth.Auth, controller.MessageList)
	return r

}
