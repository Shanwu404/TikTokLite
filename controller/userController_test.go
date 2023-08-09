package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func SendRequest(method string, url string, Body io.Reader) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, Body)
	if err != nil {
		log.Fatal("NewRequest failed", err)
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do failed", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ReadAll failed", err)
		return
	}
	fmt.Println(string(body))
}

func TestRegister(t *testing.T) {
	// 用户注册——成功
	url1 := "http://localhost:8080/douyin/user/register/?username=lux1&password=123123"
	method1 := "POST"
	SendRequest(method1, url1, nil)

	// 用户注册——用户名已存在
	url2 := "http://localhost:8080/douyin/user/register/?username=lux&password=123123"
	method2 := "POST"
	SendRequest(method2, url2, nil)

	// // 用户注册——用户名为空
	// url3 := "http://localhost:8080/douyin/user/register/?username=&password=123123"
	// method3 := "POST"
	// SendRequest(method3, url3, nil)
}

func TestLogin(t *testing.T) {
	// 用户登录——成功
	url1 := "http://localhost:8080/douyin/user/login/?username=lux&password=123123"
	method1 := "POST"
	SendRequest(method1, url1, nil)
	// // 用户登录——密码错误
	url2 := "http://localhost:8080/douyin/user/login/?username=lux&password=122"
	method2 := "POST"
	SendRequest(method2, url2, nil)
	// 用户登录——用户名不存在
	url3 := "http://localhost:8080/douyin/user/login/?username=qqly&password=122"
	method3 := "POST"
	SendRequest(method3, url3, nil)
}

func TestGetUserInfo(t *testing.T) {
	// 获取用户信息——成功
	url1 := "http://localhost:8080/douyin/user/?user_id=1"
	method1 := "GET"
	SendRequest(method1, url1, nil)
	// 获取用户信息——失败
	url2 := "http://localhost:8080/douyin/user/?user_id=100"
	method2 := "GET"
	SendRequest(method2, url2, nil)
}
