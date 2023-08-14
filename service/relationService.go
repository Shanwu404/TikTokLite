package service

import (
	"github.com/Shanwu404/TikTokLite/dao"
)

type RelationService interface {
	// Follow userId关注followId用户
	Follow(userId int64, followId int64) error

	// UnFollow userId取关followId用户
	UnFollow(userId int64, followId int64) error

	// IsFollowed 查询id1是否关注id2用户
	JudgeIsFollowById(id1 int64, id2 int64) bool

	// CountFollowers 获取用户粉丝数
	CountFollowers(id int64) int64

	// CountFollowings 获取用户关注数
	CountFollowees(id int64) int64

	// GetFollowList 获取用户关注列表
	GetFollowList(userId int64) ([]dao.UserResp, error)

	// GetFollowerList 获取用户粉丝列表
	GetFollowerList(userId int64) ([]dao.UserResp, error)

	// GetFriendList 获取用户好友列表
	//GetFriendList(userId int64) ([]dao.UserResp, error)
}
