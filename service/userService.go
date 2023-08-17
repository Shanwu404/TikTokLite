package service

import "github.com/Shanwu404/TikTokLite/dao"

type UserService interface {
	// QueryUserByUsername 根据name获取User对象
	QueryUserByUsername(username string) (dao.User, error)

	// QueryUserByID 根据id获取User对象
	QueryUserByID(id int64) (dao.User, error)

	// Register 用户注册
	Register(username string, password string) (int64, int32, string)

	// Login 用户登录
	Login(username string, password string) (int32, string)

	// IsUserIdExist 查询用户ID是否存在
	IsUserIdExist(id int64) bool
}
