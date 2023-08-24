package service

import (
	"fmt"
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
)

type RelationServiceImpl struct {
}

func NewRelationService() RelationService {
	return &RelationServiceImpl{}
}

func (rs *RelationServiceImpl) Follow(userId int64, followId int64) (bool, error) {
	// 检查用户是否存在
	usi := NewUserService()
	if isExisted := usi.IsUserIdExist(followId); !isExisted {
		return false, fmt.Errorf("user %d does not exist", followId)
	} else if isExisted = usi.IsUserIdExist(userId); !isExisted {
		return false, fmt.Errorf("user %d does not exist", userId)
	}

	// 不能关注自己
	if userId == followId {
		return false, fmt.Errorf("can not follow yourself")
	}

	// 检查是否已经关注了followId
	isFollowed, err := rs.IsFollowed(userId, followId)
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

	// 将新关注关系添加到Redis缓存
	redisKey := fmt.Sprintf("relation:follow:%d", userId)
	redis.RDb.SAdd(redis.Ctx, redisKey, followId)
	redis.RDb.Expire(redis.Ctx, redisKey, 2*time.Hour)

	return true, nil
}

func (rs *RelationServiceImpl) UnFollow(userId int64, followId int64) (bool, error) {
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

	// 从Redis中移除关注关系
	redisKey := fmt.Sprintf("relation:follow:%d", userId)
	redis.RDb.SRem(redis.Ctx, redisKey, followId)

	return true, nil
}

func (rs *RelationServiceImpl) IsFollowed(userId int64, followId int64) (bool, error) {
	// 从Redis中查询关注关系
	redisKey := fmt.Sprintf("relation:follow:%d", userId)
	isFollowed, err := redis.RDb.SIsMember(redis.Ctx, redisKey, followId).Result()

	if err == nil && isFollowed {
		// 说明Redis中存在该关注关系
		return true, nil
	}

	// 如果Redis中没有关注关系，则查询数据库
	isFollowed, err = dao.IsFollowed(userId, followId)
	if err != nil {
		return false, err
	}

	// 如果数据库中存在关注关系，则将其存入Redis缓存
	if isFollowed {
		redis.RDb.SAdd(redis.Ctx, redisKey, followId)
		redis.RDb.Expire(redis.Ctx, redisKey, 2*time.Hour)
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
