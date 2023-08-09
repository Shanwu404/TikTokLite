package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Shanwu404/TikTokLite/controller"
	"github.com/Shanwu404/TikTokLite/controller/auth"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// public directory is used to serve static resources
	// r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", auth.Auth, controller.Feed)
	apiRouter.GET("/user/", auth.Auth, controller.UserInfo)

	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", auth.Auth, controller.Publish)
	apiRouter.GET("/publish/list/", auth.Auth, controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", auth.Auth, controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", auth.Auth, controller.FavoriteList)
	apiRouter.POST("/comment/action/", auth.Auth, controller.CommentAction)
	apiRouter.GET("/comment/list/", auth.Auth, controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", auth.Auth, controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", auth.Auth, controller.FollowList)
	apiRouter.GET("/relation/follower/list/", auth.Auth, controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", auth.Auth, controller.FriendList)
	apiRouter.GET("/message/chat/", auth.Auth, controller.ChatRecord)
	apiRouter.POST("/message/action/", auth.Auth, controller.MessageAction)
	return r
}
