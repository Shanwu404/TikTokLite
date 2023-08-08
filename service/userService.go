package service

import "github.com/Shanwu404/TikTokLite/dao"

type UserService interface {
	// QueryUserByUsername 根据name获取User对象
	QueryUserByUsername(username string) (dao.User, error)

	// QueryUserByID 根据id获取User对象
	QueryUserByID(id uint64) (dao.User, error)

	// QueryUserRespByID 根据id获取UserResp对象
	QueryUserRespByID(id uint64) (dao.UserResp, error)

	// Register 用户注册
	Register(username string, password string) (uint64, int32, string)

	// Login 用户登录
	Login(username string, password string) (int32, string)
}
