package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/middleware/rabbitmq"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
)

type RelationServiceImpl struct {
}

func NewRelationService() RelationService {
	return &RelationServiceImpl{}
}

// Follow 关注用户
func (rs *RelationServiceImpl) Follow(userId int64, followId int64) (bool, error) {
	// 检查用户是否存在
	usi := NewUserService()
	if isExisted := usi.IsUserIdExist(followId); !isExisted {
		logger.Errorln("user", followId, "does not exist")
		return false, fmt.Errorf("user %d does not exist", followId)
	} else if isExisted = usi.IsUserIdExist(userId); !isExisted {
		logger.Errorln("user", userId, "does not exist")
		return false, fmt.Errorf("user %d does not exist", userId)
	}

	// 不能关注自己
	if userId == followId {
		logger.Errorln("can not follow yourself")
		return false, fmt.Errorf("can not follow yourself")
	}

	// 检查是否已经关注了followId
	isFollowed, err := rs.IsFollowed(userId, followId)
	if err != nil {
		return false, err
	}
	if isFollowed {
		logger.Errorln("userId", userId, "has already followed user", followId)
		return false, fmt.Errorf("userId %d has already followed user %d", userId, followId)
	}

	// 启用消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.FormatInt(userId, 10))
	sb.WriteString(",")
	sb.WriteString(strconv.FormatInt(followId, 10))
	rabbitmq.RabbitMQRelationAdd.Producer(sb.String())

	// 保证数据一致性：主动使count缓存失效
	redisFollowCntKey := utils.RelationFollowCntKey + strconv.FormatInt(userId, 10)
	redisFollowerCntKey := utils.RelationFollowerCntKey + strconv.FormatInt(followId, 10)
	redis.RDb.Del(redis.Ctx, redisFollowCntKey, redisFollowerCntKey)

	return true, nil
}

// UnFollow 取消关注
func (rs *RelationServiceImpl) UnFollow(userId int64, followId int64) (bool, error) {
	// 检查用户是否已经关注了followId
	isFollowed, err := dao.IsFollowed(userId, followId)
	if err != nil {
		return false, err
	}
	if !isFollowed {
		logger.Errorln("userId", userId, "has not followed user", followId)
		return false, fmt.Errorf("userId %d has not followed user %d", userId, followId)
	}

	// 启用消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.FormatInt(userId, 10))
	sb.WriteString(",")
	sb.WriteString(strconv.FormatInt(followId, 10))
	rabbitmq.RabbitMQRelationDel.Producer(sb.String())

	// 保证数据一致性：主动使count缓存失效
	redisFollowCntKey := utils.RelationFollowCntKey + strconv.FormatInt(userId, 10)
	redisFollowerCntKey := utils.RelationFollowerCntKey + strconv.FormatInt(followId, 10)
	redis.RDb.Del(redis.Ctx, redisFollowCntKey, redisFollowerCntKey)

	return true, nil
}

// IsFollowed 检查是否已经关注了followId
func (rs *RelationServiceImpl) IsFollowed(userId int64, followId int64) (bool, error) {
	// 从Redis中查询关注关系
	redisFollowKey := utils.RelationFollowKey + strconv.FormatInt(userId, 10)
	isFollowed, err := redis.RDb.SIsMember(redis.Ctx, redisFollowKey, followId).Result()
	if err == nil && isFollowed {
		// 说明Redis中存在该关注关系, 更新过期时间
		redis.RDb.Expire(redis.Ctx, redisFollowKey, utils.RelationFollowKeyTTL)
		return true, nil
	}

	// 如果Redis中没有关注关系，则查询数据库
	isFollowed, err = dao.IsFollowed(userId, followId)
	if err != nil {
		return false, err
	}
	logger.Infof("user %d has followed user %d: %t", userId, followId, isFollowed)

	// 如果数据库中存在关注关系，则将其存入Redis缓存
	if isFollowed {
		redis.RDb.SAdd(redis.Ctx, redisFollowKey, followId)
		redis.RDb.Expire(redis.Ctx, redisFollowKey, utils.RelationFollowKeyTTL)
	}

	return isFollowed, nil
}

// CountFollowers 计算用户粉丝数
func (rs *RelationServiceImpl) CountFollowers(userId int64) (int64, error) {
	// 从Redis中获取用户粉丝数
	redisFollowerCntKey := utils.RelationFollowerCntKey + strconv.FormatInt(userId, 10)
	followerCnt, err := redis.RDb.Get(redis.Ctx, redisFollowerCntKey).Int64()
	if err == nil {
		// 说明Redis中存在该用户粉丝数, 更新过期时间
		redis.RDb.Expire(redis.Ctx, redisFollowerCntKey, utils.RelationFollowerCntKeyTTL)
		return followerCnt, nil
	}

	// 如果Redis中没有用户粉丝数，则从数据库中获取
	followerCnt, err = dao.CountFollowers(userId)
	if err != nil {
		return 0, err
	}
	logger.Infof("user %d has %d followers", userId, followerCnt)

	// 将用户粉丝数存入Redis
	redis.RDb.Set(redis.Ctx, redisFollowerCntKey, followerCnt, utils.RelationFollowerCntKeyTTL)

	return followerCnt, nil
}

// CountFollows 计算用户关注数
func (rs *RelationServiceImpl) CountFollows(userId int64) (int64, error) {
	// 从Redis中获取用户关注数
	redisFollowCntKey := utils.RelationFollowCntKey + strconv.FormatInt(userId, 10)
	followCnt, err := redis.RDb.Get(redis.Ctx, redisFollowCntKey).Int64()
	if err == nil {
		// 说明Redis中存在该用户关注数, 更新过期时间
		redis.RDb.Expire(redis.Ctx, redisFollowCntKey, utils.RelationFollowCntKeyTTL)
		return followCnt, nil
	}

	// 如果Redis中没有用户关注数，则从数据库中获取
	followCnt, err = dao.CountFollows(userId)
	if err != nil {
		return 0, err
	}
	logger.Infof("user %d has %d follows", userId, followCnt)

	// 将用户关注数存入Redis
	redis.RDb.Set(redis.Ctx, redisFollowCntKey, followCnt, utils.RelationFollowCntKeyTTL)

	return followCnt, nil
}

// GetFollowList 获取用户关注ID列表
func (rs *RelationServiceImpl) GetFollowList(userId int64) ([]int64, error) {
	// 获取用户关注ID列表
	followId, err := dao.QueryFollowsIdByUserId(userId)
	if err != nil {
		return nil, err
	}
	logger.Infoln("get follow list success")
	return followId, nil
}

// GetFollowerList 获取用户粉丝ID列表
func (rs *RelationServiceImpl) GetFollowerList(userId int64) ([]int64, error) {
	// 获取用户粉丝ID列表
	followerId, err := dao.QueryFollowersIdByUserId(userId)
	if err != nil {
		return nil, err
	}
	logger.Infoln("get follower list success")
	return followerId, nil
}

// GetFriendList 获取用户好友ID列表
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
	logger.Infoln("get friend list success")
	return friendList, nil
}
