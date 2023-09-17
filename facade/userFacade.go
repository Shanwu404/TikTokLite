package facade

import (
	"github.com/Shanwu404/TikTokLite/middleware/auth"
	"github.com/Shanwu404/TikTokLite/service"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfoResponse struct {
	Response
	UserInfo service.UserInfoParams `json:"user"` // 新的userinfo结构体
}

type LoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type RegisterResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserFacade struct {
	userService service.UserService
}

func NewUserFacade() *UserFacade {
	return &UserFacade{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
func (uf *UserFacade) Register(req RegisterRequest) RegisterResponse {
	userId, code, message := uf.userService.Register(req.Username, req.Password)

	if code != 0 {
		return RegisterResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		}
	} else {
		token, _ := auth.GenerateToken(req.Username, userId)
		return RegisterResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   userId,
			Token:    token,
		}
	}
}

// Login 用户登录
func (uf *UserFacade) Login(req LoginRequest) LoginResponse {
	userId, code, message := uf.userService.Login(req.Username, req.Password)

	if code != 0 {
		return LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
		}
	} else {
		token, _ := auth.GenerateToken(req.Username, userId)
		return LoginResponse{
			Response: Response{StatusCode: code, StatusMsg: message},
			UserId:   userId,
			Token:    token,
		}
	}
}

// GetUserInfo 获取用户信息
func (uf *UserFacade) GetUserInfo(userId int64) UserInfoResponse {
	userinfo, err := uf.userService.QueryUserInfoByID(userId)
	if err != nil {
		return UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		}
	}
	return UserInfoResponse{
		Response: Response{StatusCode: 0},
		UserInfo: userinfo,
	}
}
