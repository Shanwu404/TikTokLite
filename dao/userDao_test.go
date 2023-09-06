package dao

import (
	"fmt"
	"testing"
)

func TestInsertUser(t *testing.T) {
	Init()
	user1 := User{
		Username: "John",
		Password: "123456",
	}
	userId, err := InsertUser(user1)
	fmt.Println(userId)
	fmt.Println(err)
}

func TestQueryUserByID(t *testing.T) {
	Init()
	user, err := QueryUserByID(1)
	fmt.Println(user)
	fmt.Println(err)
}

func TestQueryUserByUsername(t *testing.T) {
	Init()
	user, err := QueryUserByUsername("John")
	fmt.Println(user)
	fmt.Println(err)
}

func TestIsUserIdExist(t *testing.T) {
	Init()
	flag := IsUserIdExist(1)
	fmt.Println(flag)
}
