package controller

import (
	"net/http"
	"strconv"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type RelationController struct {
	userService     service.UserService
	relationService service.RelationService
	/* 获取UserInfo所需要的接口 */
	videoService service.VideoService
}

func NewRelationController() *RelationController {
	return &RelationController{
		userService:     service.NewUserService(),
		relationService: service.NewRelationService(),
		videoService:    service.NewVideoService(),
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
				StatusMsg:  "get follower list failed",
			},
			UserList: nil,
		})
		return
	}
	// 4. 将用户粉丝列表转换为UserInfo列表
	var userInfoList []UserInfo
	for _, followerId := range followerList {
		UserInfo := rc.completeUserInfo(followerId)
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

/*--------------------------------组装用户信息----------------------------*/
func (rc *RelationController) completeUserInfo(userId int64) UserInfo {
	user, _ := rc.userService.QueryUserByID(userId)
	followCount, _ := rc.relationService.CountFollows(userId)
	followerCount, _ := rc.relationService.CountFollowers(userId)
	workCount := int64(len(rc.videoService.GetVideoListByUserId(userId)))
	// TODO
	// TotalFavorited :=
	// FavoriteCount :=

	return UserInfo{
		Id:              user.ID,
		Username:        user.Username,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        false,
		Avatar:          "https://mary-aliyun-img.oss-cn-beijing.aliyuncs.com/typora/202308171029672.jpg",
		BackgroundImage: "https://mary-aliyun-img.oss-cn-beijing.aliyuncs.com/typora/202308171007006.jpg",
		Signature:       "TikTokLite Signature",
		TotalFavorited:  123456789,
		WorkCount:       workCount,
		FavoriteCount:   2,
	}
}
