package controller

import (
	"net/http"

	"github.com/Shanwu404/TikTokLite/middleware/auth"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/Shanwu404/TikTokLite/utils"
	"github.com/Shanwu404/TikTokLite/utils/validation"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

type UserResponse struct {
	Response
	UserInfo service.UserInfoParams `json:"user"` // 新的userinfo结构体
}

type LoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

// Register POST /douyin/user/register/ 用户注册
func (uc *UserController) Register(c *gin.Context) {
	// 获取客户端IP地址
	clientIP := c.ClientIP()
	// 检查IP是否被限制
	if isLimited := utils.IsRateLimited(clientIP); isLimited {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "注册请求次数过多，请稍后再试"},
		})
		return
	}

	// 解析注册请求参数并校验
	req, isValid := validation.RegisterParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或密码不合法"},
		})
		return
	}

	username := req.Username
	password := req.Password

	userId, code, message := uc.userService.Register(username, password)

	if code != 0 {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		})
		return
	} else {
		token, _ := auth.GenerateToken(username, userId)
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   userId,
			Token:    token,
		})
		return
	}
}

// Login POST /douyin/user/login/ 用户登录
func (uc *UserController) Login(c *gin.Context) {
	// 获取客户端IP地址
	clientIP := c.ClientIP()
	// 检查IP是否被限制
	if isLimited := utils.IsRateLimited(clientIP); isLimited {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "登录请求次数过多, 请稍后再试"},
		})
		return
	}

	// 解析登录请求参数并校验
	req, isValid := validation.LoginParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或密码不合法, 请检查"},
		})
		return
	}

	username := req.Username
	password := req.Password

	code, message := uc.userService.Login(username, password)
	if code != 0 {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		})
	} else {
		user, err := uc.userService.QueryUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, LoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Internal Server Error"},
			})
			return
		}
		token, _ := auth.GenerateToken(user.Username, user.ID)
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   user.ID,
			Token:    token,
		})
		return
	}
}

// GetUserInfo GET /douyin/user/ 用户信息
func (uc *UserController) GetUserInfo(c *gin.Context) {

	// 解析请求参数并校验
	userId, isValid := validation.GetUserInfoParseAndValidateParams(c)
	if !isValid {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户ID不合法"},
		})
		return
	}

	userinfo, err := uc.userService.QueryUserInfoByID(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		UserInfo: userinfo,
	})
}
