package main

import (
	"github.com/Shanwu404/TikTokLite/controller"
	"github.com/Shanwu404/TikTokLite/controller/auth"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	apiRouter := r.Group("/douyin")

	userService := service.NewUserService()                     // 实例化 UserService
	userController := controller.NewUserController(userService) // 实例化 UserController

	videoController := controller.NewVideoController()
	// basic apis
	apiRouter.GET("/feed/", videoController.Feed)
	apiRouter.POST("/publish/action/", auth.Auth, videoController.PublishAction)
	apiRouter.POST("/publish/list/", auth.Auth, videoController.PublishList)

	apiRouter.POST("/user/register/", userController.Register)
	apiRouter.POST("/user/login/", userController.Login)
	apiRouter.GET("/user/", userController.UserInfo)

	return r

}
