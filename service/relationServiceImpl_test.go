package service

import (
	"fmt"
	"testing"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
)

func RelationServiceInit() {
	dao.Init()
	redis.InitRedis()
}

func TestRelationServiceImpl_Follow(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	flag, err := rsi.Follow(7, 14)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_UnFollow(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	flag, err := rsi.UnFollow(3, 14)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_IsFollowed(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	flag, err := rsi.IsFollowed(3, 14)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_CountFollowers(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	flag, err := rsi.CountFollowers(1001)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_CountFollows(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	flag, err := rsi.CountFollows(7)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_GetFollowList(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	followList, err := rsi.GetFollowList(1000)
	fmt.Println(followList, err)
}

func TestRelationServiceImpl_GetFollowerList(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	followerList, err := rsi.GetFollowerList(7)
	fmt.Println(followerList, err)
}

func TestRelationServiceImpl_GetFriendList(t *testing.T) {
	RelationServiceInit()
	rsi := NewRelationService()
	friendList, err := rsi.GetFriendList(1000)
	fmt.Println(friendList, err)
}
