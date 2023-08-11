package main

import (
	"github.com/Shanwu404/TikTokLite/controller"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	apiRouter := r.Group("/douyin")

	userService := service.NewUserService()                     // 实例化 UserService
	userController := controller.NewUserController(userService) // 实例化 UserController

	// basic apis
	// apiRouter.GET("/feed/", controller.Feed)
	// apiRouter.POST("/publish/action/", auth.Auth, controller.Publish)

	apiRouter.POST("/user/register/", userController.Register)
	apiRouter.POST("/user/login/", userController.Login)
	apiRouter.GET("/user/", userController.GetUserInfo)

	return r

}
