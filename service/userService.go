package service

import "github.com/Shanwu404/TikTokLite/dao"

// 新增UserInfoParams结构体,避免暴露dao.UserInfo结构体
type UserInfoParams struct {
	Id              int64  `json:"id"`               // 用户ID
	Username        string `json:"name"`             // 用户名
	FollowCount     int64  `json:"follow_count"`     // 关注数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝数
	IsFollow        bool   `json:"is_follow"`        // 是否关注
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 背景图片
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数
	WorkCount       int64  `json:"work_count"`       // 作品数
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
}

type UserService interface {
	// QueryUserByUsername 根据name获取User对象
	QueryUserByUsername(username string) (dao.User, error)

	// QueryUserByID 根据id获取User对象
	QueryUserByID(id int64) (dao.User, error)

	// Register 用户注册
	Register(username string, password string) (int64, int32, string)

	// Login 用户登录
	Login(username string, password string) (int64, int32, string)

	// IsUserIdExist 查询用户ID是否存在
	IsUserIdExist(id int64) bool

	// QueryUserInfoByID 根据用户ID查询用户信息
	QueryUserInfoByID(id int64) (UserInfoParams, error)
}
