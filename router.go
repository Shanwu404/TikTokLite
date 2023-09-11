package main

import (
	"github.com/Shanwu404/TikTokLite/controller"
	"github.com/Shanwu404/TikTokLite/middleware/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	// pprof.Register(r) // 性能分析
	apiRouter := r.Group("/douyin")
	
	userController := controller.NewUserController()
	relationController := controller.NewRelationController()
	videoController := controller.NewVideoController()
	commentController := controller.NewCommentController()
	messageController := controller.NewMessageController()
	likeController := controller.NewLikeController()

	// basic apis
	apiRouter.GET("/feed/", auth.Auth, videoController.Feed)
	apiRouter.POST("/publish/action/", auth.Auth, videoController.PublishAction)
	apiRouter.GET("/publish/list/", auth.Auth, videoController.PublishList)

	apiRouter.POST("/user/register/", userController.Register)
	apiRouter.POST("/user/login/", userController.Login)
	apiRouter.GET("/user/", auth.Auth, userController.GetUserInfo)

	apiRouter.GET("/relation/follow/list/", auth.Auth, relationController.FollowsList)
	apiRouter.GET("/relation/follower/list/", auth.Auth, relationController.FollowersList)
	apiRouter.POST("/relation/action/", auth.Auth, relationController.RelationAction)
	apiRouter.GET("/relation/friend/list/", auth.Auth, relationController.FriendList)

	apiRouter.POST("/comment/action/", auth.Auth, commentController.CommentAction)
	apiRouter.GET("/comment/list/", auth.Auth, commentController.CommentList)

	apiRouter.POST("/message/action/", auth.Auth, messageController.MessageAction)
	apiRouter.GET("/message/chat/", auth.Auth, messageController.MessageList)

	apiRouter.POST("/favorite/action/", auth.Auth, likeController.FavoriteAction)
	apiRouter.GET("/favorite/list/", auth.Auth, likeController.FavoriteList)
	return r

}
