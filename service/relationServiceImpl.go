package service

import (
	"fmt"

	"github.com/Shanwu404/TikTokLite/dao"
)

func NewRelationService() *RelationServiceImpl {
	return &RelationServiceImpl{
		UserService: &UserServiceImpl{},
	}
}

type RelationServiceImpl struct {
	UserService *UserServiceImpl
}

func (rs *RelationServiceImpl) Follow(userId int64, followId int64) (bool, error) {
	// 检查用户是否已经关注了followId
	isFollowed, err := dao.IsFollowed(userId, followId)
	if err != nil {
		return false, err
	}
	if isFollowed {
		return false, fmt.Errorf("userId %d has already followed user %d", userId, followId)
	}

	// 插入新的关注关系
	if err := dao.InsertFollow(userId, followId); err != nil {
		return false, err
	}

	return true, nil
}

func (rs *RelationServiceImpl) Unfollow(userId int64, followId int64) (bool, error) {
	// 检查用户是否已经关注了followId
	isFollowed, err := dao.IsFollowed(userId, followId)
	if err != nil {
		return false, err
	}
	if !isFollowed {
		return false, fmt.Errorf("userId %d has not followed user %d", userId, followId)
	}

	// 删除关注关系
	if err := dao.DeleteFollow(userId, followId); err != nil {
		return false, err
	}

	return true, nil
}

func (rs *RelationServiceImpl) IsFollowed(userId int64, followId int64) (bool, error) {
	// 检查用户是否已经关注了followId
	isFollowed, err := dao.IsFollowed(userId, followId)
	if err != nil {
		return false, err
	}

	return isFollowed, nil
}

func (rs *RelationServiceImpl) CountFollowers(userId int64) (int64, error) {
	// 获取用户粉丝数
	count, err := dao.CountFollowers(userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (rs *RelationServiceImpl) CountFollows(userId int64) (int64, error) {
	// 获取用户关注数
	count, err := dao.CountFollows(userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (rs *RelationServiceImpl) GetFollowList(userId int64) ([]int64, error) {
	// 获取用户关注ID列表
	followId, err := dao.QueryFollowsIdByUserId(userId)
	if err != nil {
		return nil, err
	}

	return followId, nil
}

func (rs *RelationServiceImpl) GetFollowerList(userId int64) ([]int64, error) {
	// 获取用户粉丝ID列表
	followerId, err := dao.QueryFollowersIdByUserId(userId)
	if err != nil {
		return nil, err
	}

	return followerId, nil
}

func (rs *RelationServiceImpl) GetFriendList(userId int64) ([]int64, error) {
	// 查询用户关注ID列表
	followId, err := dao.QueryFollowsIdByUserId(userId)
	if err != nil {
		return nil, err
	}

	// 查询用户粉丝ID列表
	followerId, err := dao.QueryFollowersIdByUserId(userId)
	if err != nil {
		return nil, err
	}

	// 找出既是关注者又是被关注者的ID
	friendMap := make(map[int64]bool)
	for _, id := range followId {
		friendMap[id] = true
	}

	var friendList []int64
	for _, id := range followerId {
		if _, exists := friendMap[id]; exists {
			friendList = append(friendList, id)
		}
	}

	return friendList, nil
}
