package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Shanwu404/TikTokLite/service"
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
	UserList []UserInfo `json:"user_list"`
}

type FriendListResponse struct {
	UserListResponse
}

// RelationAction POST /douyin/relation/action/ 关注/取消关注
func (rc *RelationController) RelationAction(c *gin.Context) {

	// 1. 取出用户id
	userId := c.GetInt64("id")

	// 2. 判断to_user_id解析是否有误
	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "followId error"})
		return
	}

	// 3. 判断actionType解释是否有误
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "actionType error"})
		return
	}

	// 4. 执行关注/取关操作
	switch {
	case actionType == 1:
		// 执行关注
		flag, err := rc.relationService.Follow(userId, followId)
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

	case actionType == 2:
		// 执行取关
		flag, err := rc.relationService.UnFollow(userId, followId)
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

// FollowersList GET /douyin/relation/follow/list/ 获取粉丝列表
func (rc *RelationController) FollowersList(c *gin.Context) {
	// 1. 取出用户id
	userId := c.GetInt64("id")

	// 2. 检查查询参数中的用户ID是否存在并且与当前用户ID相匹配
	queryUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userId != queryUserId {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "user_id error",
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
				StatusMsg:  "获取粉丝列表失败",
			},
			UserList: nil,
		})
		return
	}
	// TODO 4. 将用户粉丝列表转换为UserInfo列表
	log.Println("followerList: ", followerList)
}
