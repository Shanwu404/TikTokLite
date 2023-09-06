package dao

import (
	"github.com/Shanwu404/TikTokLite/log/logger"
)

type User struct {
	ID       int64
	Username string
	Password string `json:"-"`
}

// InsertUser 新增用户
func InsertUser(user User) (int64, error) {
	err := db.Create(&user).Error
	if err != nil {
		logger.Errorln(err)
		return 0, err
	}
	return user.ID, nil
}

// QueryUserByID 根据ID查询User
func QueryUserByID(id int64) (User, error) {
	var user User
	result := db.Where("id = ?", id).First(&user)
	if err := result.Error; err != nil {
		logger.Errorln(err, id)
		return User{}, err
	}
	return user, nil
}

// QueryUserByUsername 根据Username查询User
func QueryUserByUsername(username string) (User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)

	if err := result.Error; err != nil {
		logger.Errorln(err, username)
		return User{}, err
	}
	return user, nil
}

// 查询用户ID是否存在
func IsUserIdExist(id int64) bool {
	user := &User{}
	result := db.Where("id = ?", id).First(user)
	if err := result.Error; err != nil {
		logger.Errorln(err)
		return false
	}
	return true
}
