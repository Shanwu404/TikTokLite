package controller

import (
	"net/http"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/Shanwu404/TikTokLite/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userFacade *facade.UserFacade
}

func NewUserController() *UserController {
	return &UserController{
		userFacade: facade.NewUserFacade(),
	}
}

// Register POST /douyin/user/register/ 用户注册
func (uc *UserController) Register(c *gin.Context) {
	// 获取客户端IP地址
	clientIP := c.ClientIP()
	// 检查IP是否被限制
	if isLimited := utils.IsRateLimited(clientIP); isLimited {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "注册请求次数过多，请稍后再试"})
		return
	}

	// 解析注册请求参数并校验
	req, isValid := RegisterParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户名或密码不合法"})
		return
	}

	// 执行注册操作
	RegisterResponse := uc.userFacade.Register(req)

	// 返回响应
	c.JSON(http.StatusOK, RegisterResponse)
}

// Login POST /douyin/user/login/ 用户登录
func (uc *UserController) Login(c *gin.Context) {
	// 获取客户端IP地址
	clientIP := c.ClientIP()
	// 检查IP是否被限制
	if isLimited := utils.IsRateLimited(clientIP); isLimited {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "登录请求次数过多, 请稍后再试"})
		return
	}

	// 解析登录请求参数并校验
	req, isValid := LoginParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户名或密码不合法, 请检查"})
		return
	}

	// 执行登录操作
	LoginResponse := uc.userFacade.Login(req)

	// 返回响应
	c.JSON(http.StatusOK, LoginResponse)
}

// GetUserInfo GET /douyin/user/ 用户信息
func (uc *UserController) GetUserInfo(c *gin.Context) {

	// 解析请求参数并校验
	userId, isValid := GetUserInfoParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户ID不合法"})
		return
	}

	// 执行查询操作
	UserResponse := uc.userFacade.GetUserInfo(userId)

	// 返回响应
	c.JSON(http.StatusOK, UserResponse)
}
