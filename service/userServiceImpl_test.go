package service

import (
	"fmt"
	"testing"

	"github.com/Shanwu404/TikTokLite/dao"
)

func UserServiceImplInit() {
	dao.Init()
}

func TestUserServiceImpl_QueryUserByName(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	user, err := usi.QueryUserByUsername("Lihua")
	fmt.Println(user)
	fmt.Println(err)
}

func TestUserServiceImpl_QueryUserByID(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	user, err := usi.QueryUserByID(1)
	fmt.Println(user)
	fmt.Println(err)
}

func TestUserServiceImpl_Register(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	id, code, message := usi.Register("Lqs1", "password")
	fmt.Println(id, code, message)
}

func TestUserServiceImpl_Login(t *testing.T) {
	UserServiceImplInit()
	usi := UserServiceImpl{}
	code, message := usi.Login("Lqs", "1000")
	fmt.Println(code, message)
}

func TestComparePasswords(t *testing.T) {
	hashedPassword := "$2a$10$GRozN2nx7FZncQO/Zhx2yer4vd1Xbey4VC1DtCNjPtZnpvufWVvgG"
	originalPassword := "1000"
	match := ComparePasswords(hashedPassword, originalPassword)
	if match {
		fmt.Println("密码匹配!")
	} else {
		fmt.Println("密码不匹配.")
	}
}

func TestQueryUserInfoByID(t *testing.T) {
	UserServiceImplInit()
	usi := NewUserService()
	userInfo, err := usi.QueryUserInfoByID(14)
	fmt.Println(userInfo)
	fmt.Println(err)
}
