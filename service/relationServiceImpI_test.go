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
	err := rsi.Follow(1000, 1001)
	fmt.Println(err)
}

func TestRelationServiceImpl_UnFollow(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	err := rsi.UnFollow(1000, 1001)
	fmt.Println(err)
}

func TestRelationServiceImpl_JudgeIsFollowById(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	flag := rsi.JudgeIsFollowById(1000, 1001)
	fmt.Println(flag)
}

func TestRelationServiceImpl_CountFollowers(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	cnt := rsi.CountFollowers(1000)
	fmt.Println(cnt)
}

func TestRelationServiceImpl_CountFollowings(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	cnt := rsi.CountFollowings(1000)
	fmt.Println(cnt)
}

func TestRelationServiceImpl_GetFollowerList(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	followers, err := rsi.GetFollowerList(1000)
	fmt.Println(followers)
	fmt.Println(err)
}

func TestRelationServiceImpl_GetFollowList(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	follows, err := rsi.GetFollowList(1000)
	fmt.Println(follows)
	fmt.Println(err)
}

func TestRelationServiceImpl_GetFriendList(t *testing.T) {
	RelationServiceInit()
	rsi := RelationServiceImpl{}
	friends, err := rsi.GetFriendList(1000)
	fmt.Println(friends)
	fmt.Println(err)
}
