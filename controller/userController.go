package controller

import (
	"net/http"
	"strconv"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Response
	User dao.UserResp `json:"user"`
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

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
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
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   userId,
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
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   user.ID,
		})
		return
	}
}

// UserInfo GET /douyin/user/ 用户信息
func (uc *UserController) UserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	id, _ := strconv.ParseInt(userId, 10, 64) // 字符串转int64
	userInfo, err := uc.userService.QueryUserRespByID(id)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User does not exist"},
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     userInfo,
	})
}
