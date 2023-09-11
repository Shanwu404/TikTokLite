package controller

import (
	"net/http"
	"strconv"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/Shanwu404/TikTokLite/utils/validation"
	"github.com/gin-gonic/gin"
)

type RelationController struct {
	userService     service.UserService
	relationService service.RelationService
}

func NewRelationController() *RelationController {
	return &RelationController{
		userService:     service.NewUserService(),
		relationService: service.NewRelationService(),
	}
}

type UserListResponse struct {
	Response
	UserList []service.UserInfoParams `json:"user_list"`
}

// RelationAction POST /douyin/relation/action/ 关注/取消关注
func (rc *RelationController) RelationAction(c *gin.Context) {

	// 1. 解析关注/取关请求参数并校验
	req, isValid := validation.RelationActionParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "invalid params",
		})
		return
	}

	// 执行关注/取关操作
	switch {
	case req.ActionType == 1:
		// 执行关注
		flag, err := rc.relationService.Follow(req.UserId, req.ToUserId)
		if err != nil || !flag {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "follow success!"})
		return

	case req.ActionType == 2:
		// 执行取关
		flag, err := rc.relationService.UnFollow(req.UserId, req.ToUserId)
		if err != nil || !flag {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				StatusMsg:  err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "unfollow success!"})
		return
	}
}

// FriendsList GET /douyin/relation/follow/list/ 获取关注列表
func (rc *RelationController) FollowsList(c *gin.Context) {
	// 1. 取出用户id
	userId := c.GetInt64("id")

	// 2. 检查查询参数中的用户ID是否存在并且与当前用户ID相匹配
	queryUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userId != queryUserId {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "由于该用户隐私设置, 关注列表不可见",
			},
			UserList: nil,
		})
		return
	}

	// 3. 获取用户关注列表
	followList, err := rc.relationService.GetFollowList(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "get follow list failed",
			},
			UserList: nil,
		})
		return
	}

	// 4. 将用户关注列表转换为UserInfo列表
	var userInfoList []service.UserInfoParams
	for _, followId := range followList {
		UserInfo, _ := rc.userService.QueryUserInfoByID(followId)
		UserInfo.IsFollow, _ = rc.relationService.IsFollowed(userId, followId)
		userInfoList = append(userInfoList, UserInfo)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get follow list success",
		},
		UserList: userInfoList,
	})
}

// FollowersList GET /douyin/relation/follower/list/ 获取粉丝列表
func (rc *RelationController) FollowersList(c *gin.Context) {
	// 1. 取出用户id
	userId := c.GetInt64("id")

	// 2. 检查查询参数中的用户ID是否存在并且与当前用户ID相匹配
	queryUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userId != queryUserId {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "由于该用户隐私设置, 粉丝列表不可见",
			},
			UserList: nil,
		})
		return
	}

	// 3. 获取用户粉丝列表
	followerList, err := rc.relationService.GetFollowerList(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "get follower list failed",
			},
			UserList: nil,
		})
		return
	}

	// 4. 将用户粉丝列表转换为UserInfo列表
	var userInfoList []service.UserInfoParams
	for _, followerId := range followerList {
		UserInfo, _ := rc.userService.QueryUserInfoByID(followerId)
		UserInfo.IsFollow, _ = rc.relationService.IsFollowed(userId, followerId)
		userInfoList = append(userInfoList, UserInfo)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get follower list success",
		},
		UserList: userInfoList,
	})
}

// FriendList GET /douyin/relation/friend/list/ 获取好友列表
func (rc *RelationController) FriendList(c *gin.Context) {

	// 1. 取出用户id
	userId := c.GetInt64("id")

	// 2. 检查查询参数中的用户ID是否存在并且与当前用户ID相匹配
	queryUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userId != queryUserId {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "token鉴权失败, 好友列表不可见",
			},
			UserList: nil,
		})
		return
	}

	// 3. 获取用户好友列表
	friendList, err := rc.relationService.GetFriendList(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "get friend list failed",
			},
			UserList: nil,
		})
		return
	}

	// 4. 将用户好友列表转换为UserInfo列表
	var userInfoList []service.UserInfoParams
	for _, friendId := range friendList {
		UserInfo, _ := rc.userService.QueryUserInfoByID(friendId)
		UserInfo.IsFollow = true // 好友列表中的用户一定是关注了当前用户的
		userInfoList = append(userInfoList, UserInfo)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get friend list success",
		},
		UserList: userInfoList,
	})
}
