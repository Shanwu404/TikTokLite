package service

type RelationService interface {
	// Follow 关注followId用户
	Follow(userId int64, followId int64) (bool, error)

	// UnFollow 取消关注followId用户
	UnFollow(userId int64, followId int64) (bool, error)

	// IsFollowed 查询是否已关注followId用户
	IsFollowed(userId int64, followId int64) (bool, error)

	// CountFollowers 获取用户粉丝数
	CountFollowers(userId int64) (int64, error)

	// CountFollows 获取用户关注数
	CountFollows(userId int64) (int64, error)

	// GetFollowList 获取用户关注列表
	GetFollowList(userId int64) ([]int64, error)

	// GetFollowerList 获取用户粉丝列表
	GetFollowerList(userId int64) ([]int64, error)

	// GetFriendList 获取用户好友列表
	GetFriendList(userId int64) ([]int64, error)
}
