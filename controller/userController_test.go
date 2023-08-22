package controller

import (
	"strings"
	"testing"

	"github.com/Shanwu404/TikTokLite/middleware/auth"
)

func TestRegister(t *testing.T) {
	// 用户注册——成功
	url1 := "http://localhost:8080/douyin/user/register/?username=chy&password=123456"
	method1 := "POST"
	SendRequest(method1, url1, nil)

	// // 用户注册——用户名已存在
	// url2 := "http://localhost:8080/douyin/user/register/?username=lux&password=123123"
	// method2 := "POST"
	// SendRequest(method2, url2, nil)

	// // 用户注册——用户名为空
	// url3 := "http://localhost:8080/douyin/user/register/?username=&password=123123"
	// method3 := "POST"
	// SendRequest(method3, url3, nil)
}

func TestLogin(t *testing.T) {
	// 用户登录——成功
	url1 := "http://localhost:8080/douyin/user/login/?username=chy&password=123456"
	method1 := "POST"
	SendRequest(method1, url1, nil)
	// // // 用户登录——密码错误
	// url2 := "http://localhost:8080/douyin/user/login/?username=lux&password=122"
	// method2 := "POST"
	// SendRequest(method2, url2, nil)
	// // 用户登录——用户名不存在
	// url3 := "http://localhost:8080/douyin/user/login/?username=qqly&password=122"
	// method3 := "POST"
	// SendRequest(method3, url3, nil)
}

func TestGetUserInfo(t *testing.T) {
	token, _ := auth.GenerateToken("lux1", 7)

	// 获取用户信息——成功
	url2 := "http://localhost:8080/douyin/user/?user_id=7&token=" + token
	method2 := "GET"
	SendRequest(method2, url2, strings.NewReader(token))

	// // 获取用户信息——失败
	// url2 := "http://localhost:8080/douyin/user/?user_id=100&token=" + token
	// method2 := "GET"
	// SendRequest(method2, url2, nil)

	// 获取用户信息——鉴权失败
	// url3 := "http://localhost:8080/douyin/user/?user_id=6"
	// method3 := "GET"
	// SendRequest(method3, url3, strings.NewReader("token=123"))
}
