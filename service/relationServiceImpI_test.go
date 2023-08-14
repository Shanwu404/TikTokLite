package service

import (
	"fmt"
	"testing"

	"github.com/Shanwu404/TikTokLite/dao"
)

func RelationServiceInit() {
	dao.Init()
}

func TestRelationServiceImpl_Follow(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	flag, err := rsi.Follow(1000, 2300)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_Unfollow(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	flag, err := rsi.Unfollow(1000, 2300)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_IsFollowed(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	flag, err := rsi.IsFollowed(1000, 2300)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_CountFollowers(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	flag, err := rsi.CountFollowers(1001)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_CountFollows(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	flag, err := rsi.CountFollows(1001)
	fmt.Println(flag, err)
}

func TestRelationServiceImpl_GetFollowList(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	followList, err := rsi.GetFollowList(1000)
	fmt.Println(followList, err)
}

func TestRelationServiceImpl_GetFollowerList(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	followerList, err := rsi.GetFollowerList(1000)
	fmt.Println(followerList, err)
}

func TestRelationServiceImpl_GetFriendList(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	friendList, err := rsi.GetFriendList(1000)
	fmt.Println(friendList, err)
}
