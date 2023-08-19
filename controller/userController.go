package controller

import (
	"net/http"
	"strconv"

	"github.com/Shanwu404/TikTokLite/middleware/auth"
	"github.com/Shanwu404/TikTokLite/service"

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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register POST /douyin/user/register/ 用户注册
func (uc *UserController) Register(c *gin.Context) {
	// 可以修改为查询JSON
	username := c.Query("username")
	password := c.Query("password")

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
	username := c.Query("username")
	password := c.Query("password")

	code, message := uc.userService.Login(username, password)
	if code != 0 {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		})
	} else {
		user, err := uc.userService.QueryUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving user information"})
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
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64) // 字符串转int64
	if err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Invalid user ID format"},
		})
		return
	}

	uc.userService.IsUserIdExist(userId)
	if isExisted := uc.userService.IsUserIdExist(userId); !isExisted {
		{
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User does not exist"},
			})
			return
		}
	}

	userinfo, _ := uc.userService.QueryUserInfoByID(userId)
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		UserInfo: userinfo,
	})
}
