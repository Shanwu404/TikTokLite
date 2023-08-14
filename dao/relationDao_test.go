package dao

import (
	"fmt"
	"testing"
)

// TestInsertFollow 测试增加关注关系
func TestInsertFollow(t *testing.T) {
	Init()
	err := InsertFollow(1000, 1001)
	fmt.Println(err)
}

func TestDeleteFollow(t *testing.T) {
	Init()
	err := DeleteFollow(1000, 1001)
	fmt.Println(err)
}

func TestIsFollowed(t *testing.T) {
	Init()
	isFollow, err := IsFollowed(1000, 1001)
	fmt.Println(isFollow)
	fmt.Println(err)
}

func TestQueryFollowersIdByUserId(t *testing.T) {
	Init()
	followersId, err := QueryFollowersIdByUserId(1001)
	fmt.Println(followersId)
	fmt.Println(err)
}

func TestFollowsIdByUserId(t *testing.T) {
	Init()
	followsId, err := QueryFollowsIdByUserId(1000)
	fmt.Println(followsId)
	fmt.Println(err)
}

func TestCountFollowers(t *testing.T) {
	Init()
	count, err := CountFollowers(1001)
	fmt.Println(count)
	fmt.Println(err)
}

func TestCountFollows(t *testing.T) {
	Init()
	count, err := CountFollows(1000)
	fmt.Println(count)
	fmt.Println(err)
}
