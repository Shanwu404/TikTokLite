package controller

import (
	"net/http"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
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
func (userController *UserController) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, code, message := userController.userService.Register(request.Username, request.Password)

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
func (userController *UserController) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, message := userController.userService.Login(request.Username, request.Password)
	if code != 0 {
		c.JSON(http.StatusOK, LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		})
		return
	} else {
		user, err := userController.userService.QueryUserByUsername(request.Username)
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
