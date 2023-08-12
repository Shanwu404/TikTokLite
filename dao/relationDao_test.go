package dao

import (
	"fmt"
	"testing"
)

func TestInsertFollow(t *testing.T) {
	Init()
	err := InsertFollow(2000, 1001)
	fmt.Println(err)
}

func TestDeleteFollow(t *testing.T) {
	Init()
	err := DeleteFollow(1000, 1001)
	fmt.Println(err)
}

func TestJudgeIsFollowById(t *testing.T) {
	Init()
	err := JudgeIsFollowById(1000, 1001)
	fmt.Println(err)
}

func TestQueryFollowsIdByUserId(t *testing.T) {
	Init()
	follows, err := QueryFollowsIdByUserId(1000)
	fmt.Println(follows)
	fmt.Println(err)
}

func TestQueryFollowersIdByUserId(t *testing.T) {
	Init()
	followers, err := QueryFollowersIdByUserId(1000)
	fmt.Println(followers)
	fmt.Println(err)
}

func TestCountFollowers(t *testing.T) {
	Init()
	cnt := CountFollowers(1000)
	fmt.Println(cnt)
}

func TestCountFollowees(t *testing.T) {
	Init()
	cnt := CountFollowees(1000)
	fmt.Println(cnt)
}
