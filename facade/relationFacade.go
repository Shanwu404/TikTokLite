package facade

import (
	"sync"

	"github.com/Shanwu404/TikTokLite/service"
)

type RelationActionRequest struct {
	UserId     int64
	ToUserId   int64
	ActionType int64
}

type RelationResponse struct {
	Response
}

type UserListResponse struct {
	Response
	UserList []service.UserInfoParams `json:"user_list"`
}

type RelationFacade struct {
	relationService service.RelationService
	userService     service.UserService
}

func NewRelationFacade() *RelationFacade {
	return &RelationFacade{
		relationService: service.NewRelationService(),
		userService:     service.NewUserService(),
	}
}

// RelationAction 关注/取关
func (rf *RelationFacade) RelationAction(req RelationActionRequest) RelationResponse {
	switch {
	case req.ActionType == 1:
		// 执行关注
		flag, err := rf.relationService.Follow(req.UserId, req.ToUserId)
		if err != nil || !flag {
			return RelationResponse{
				Response: Response{StatusCode: -1, StatusMsg: err.Error()},
			}
		}
		return RelationResponse{
			Response: Response{StatusCode: 0, StatusMsg: "follow success!"},
		}

	case req.ActionType == 2:
		// 执行取关
		flag, err := rf.relationService.UnFollow(req.UserId, req.ToUserId)
		if err != nil || !flag {
			return RelationResponse{
				Response: Response{StatusCode: -1, StatusMsg: err.Error()},
			}
		}
		return RelationResponse{
			Response: Response{StatusCode: 0, StatusMsg: "unfollow success!"},
		}

	default:
		return RelationResponse{
			Response: Response{StatusCode: -1, StatusMsg: "invalid action type"},
		}
	}
}

// FollowList 获取关注列表
func (rf *RelationFacade) FollowsList(userId int64) UserListResponse {
	// 获取用户关注列表
	followList, err := rf.relationService.GetFollowList(userId)
	if err != nil {
		return UserListResponse{
			Response: Response{StatusCode: -1, StatusMsg: "get follow list failed"},
			UserList: nil,
		}
	}

	// 将用户关注列表转换为UserInfo列表 (并发处理)
	var wg sync.WaitGroup
	wg.Add(len(followList))
	userInfoList := make([]service.UserInfoParams, 0, len(followList))
	for _, followId := range followList {
		go func(followId int64) {
			defer wg.Done()
			UserInfo, _ := rf.userService.QueryUserInfoByID(followId)
			UserInfo.IsFollow, _ = rf.relationService.IsFollowed(userId, followId)
			userInfoList = append(userInfoList, UserInfo)
		}(followId)
	}
	wg.Wait()

	return UserListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "get follow list success"},
		UserList: userInfoList,
	}
}

// FollowerList 获取粉丝列表
func (rf *RelationFacade) FollowersList(userId int64) UserListResponse {
	// 获取用户粉丝列表
	followerList, err := rf.relationService.GetFollowerList(userId)
	if err != nil {
		return UserListResponse{
			Response: Response{StatusCode: -1, StatusMsg: "get follower list failed"},
			UserList: nil,
		}
	}

	// 将用户粉丝列表转换为UserInfo列表 (并发处理)
	var wg sync.WaitGroup
	wg.Add(len(followerList))
	userInfoList := make([]service.UserInfoParams, 0, len(followerList))
	for _, followerId := range followerList {
		go func(followerId int64) {
			defer wg.Done()
			UserInfo, _ := rf.userService.QueryUserInfoByID(followerId)
			UserInfo.IsFollow, _ = rf.relationService.IsFollowed(userId, followerId)
			userInfoList = append(userInfoList, UserInfo)
		}(followerId)
	}
	wg.Wait()

	return UserListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "get follower list success"},
		UserList: userInfoList,
	}
}

// FriendList 获取好友列表
func (rf *RelationFacade) FriendsList(userId int64) UserListResponse {
	// 获取用户好友列表
	friendList, err := rf.relationService.GetFriendList(userId)
	if err != nil {
		return UserListResponse{
			Response: Response{StatusCode: -1, StatusMsg: "get friend list failed"},
			UserList: nil,
		}
	}

	// 将用户好友列表转换为UserInfo列表 (并发处理)
	var wg sync.WaitGroup
	wg.Add(len(friendList))
	userInfoList := make([]service.UserInfoParams, 0, len(friendList))
	for _, friendId := range friendList {
		go func(friendId int64) {
			defer wg.Done()
			UserInfo, _ := rf.userService.QueryUserInfoByID(friendId)
			UserInfo.IsFollow, _ = rf.relationService.IsFollowed(userId, friendId)
			userInfoList = append(userInfoList, UserInfo)
		}(friendId)
	}
	wg.Wait()

	return UserListResponse{
		Response: Response{StatusCode: 0, StatusMsg: "get friend list success"},
		UserList: userInfoList,
	}
}
