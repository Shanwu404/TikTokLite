package controller

import (
	"log"
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
	UserInfo UserInfo `json:"user_info"`
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
		token = "token=" + token
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
		token = "token=" + token
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
	userId := c.Query("user_id")
	id, _ := strconv.ParseInt(userId, 10, 64) // 字符串转int64
	user, err := uc.userService.QueryUserByID(id)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User does not exist"},
		})
		return
	}
	currentUserId := c.GetInt64("id")
	log.Println("currentUserId: ", currentUserId)

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		UserInfo: UserInfo{
			Id:              user.ID,       // 用户ID
			Username:        user.Username, // 用户名
			FollowCount:     0,             // TODO: 关注数接口实现
			FollowerCount:   0,             // TODO: 粉丝数接口实现
			IsFollow:        false,         // TODO: 是否关注接口实现
			Avatar:          "",            // TODO: 头像接口实现
			BackgroundImage: "",            // TODO: 背景图片接口实现
			Signature:       "",            // TODO: 个人简介接口实现
			TotalFavorited:  0,             // TODO: 获赞数接口实现
			WorkCount:       0,             // TODO: 作品数接口实现
			FavoriteCount:   0,             // TODO: 喜欢数接口实现
		},
	})
}
