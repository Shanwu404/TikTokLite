package controller

import (
	"net/http"
	"strconv"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type RelationController struct {
	relationFacade *facade.RelationFacade
}

func NewRelationController() *RelationController {
	return &RelationController{
		relationFacade: facade.NewRelationFacade(),
	}
}

type UserListResponse struct {
	Response
	UserList []service.UserInfoParams `json:"user_list"`
}

// RelationAction POST /douyin/relation/action/ 关注/取消关注
func (rc *RelationController) RelationAction(c *gin.Context) {

	// 1. 解析关注/取关请求参数并校验
	req, isValid := RelationActionParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "invalid params",
		})
		return
	}

	// 2. 执行关注/取关操作
	RelationResponse := rc.relationFacade.RelationAction(req)

	// 3. 返回响应
	c.JSON(http.StatusOK, RelationResponse)
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
	UserListResponse := rc.relationFacade.FollowsList(userId)

	// 4. 返回响应
	c.JSON(http.StatusOK, UserListResponse)
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
	UserListResponse := rc.relationFacade.FollowersList(userId)

	// 4. 返回响应
	c.JSON(http.StatusOK, UserListResponse)
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
	UserListResponse := rc.relationFacade.FriendsList(userId)

	// 4. 返回响应
	c.JSON(http.StatusOK, UserListResponse)
}
