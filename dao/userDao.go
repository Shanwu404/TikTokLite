package dao

import (
	"log"
)

type User struct {
	ID       uint64
	Username string
	Password string
}

// InsertUser 新增用户
func InsertUser(user *User) error {
	err := Db.Create(user).Error
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

// QueryUserByID 根据ID查询User
func QueryUserByID(id uint64) (*User, error) {
	user := &User{}
	result := Db.Where("id = ?", id).First(user)
	if err := result.Error; err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return user, nil
}

// QueryUserByUsername 根据ID查询User
func QueryUserByUsername(username string) (*User, error) {
	user := &User{}
	result := Db.Where("username = ?", username).First(user)

	if err := result.Error; err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return user, nil
}

func QueryAllNames() []string {
	usernames := make([]string, 0)
	result := Db.Table("users").Pluck("username", &usernames)

	if err := result.Error; err != nil {
		log.Println(err.Error())
		return []string{}
	}

	return usernames
}
