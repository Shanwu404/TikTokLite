package main

import (
	"github.com/Shanwu404/TikTokLite/controller"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/Shanwu404/TikTokLite/utils/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	apiRouter := r.Group("/douyin")
	apiRouter.Static("tiktok", "./videos") //测试用。配置静态资源路径

	userService := service.NewUserService()                     // 实例化 UserService
	userController := controller.NewUserController(userService) // 实例化 UserController

	videoController := controller.NewVideoController()
	// basic apis
	apiRouter.GET("/feed/", videoController.Feed)
	apiRouter.POST("/publish/action/", auth.Auth, videoController.PublishAction)
	apiRouter.POST("/publish/list/", auth.Auth, videoController.PublishList)

	apiRouter.POST("/user/register/", userController.Register)
	apiRouter.POST("/user/login/", userController.Login)
	apiRouter.GET("/user/", userController.GetUserInfo)

	apiRouter.POST("/comment/action/", auth.Auth, controller.CommentAction)
	apiRouter.GET("/comment/list/", auth.Auth, controller.CommentList)

	apiRouter.POST("/message/action/", auth.Auth, controller.MessageAction)
	apiRouter.GET("/message/chat/", auth.Auth, controller.MessageList)
	return r

}
