package controller

import (
	"strconv"
	"unicode"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/gin-gonic/gin"
)

func LoginParseAndValidateParams(c *gin.Context) (facade.LoginRequest, bool) {
	req := facade.LoginRequest{
		Username: c.Query("username"),
		Password: c.Query("password"),
	}

	// 检查用户名是否合法
	if isValid := IsValidUsername(req.Username); !isValid {
		logger.Infoln("Invalid username:", req.Username)
		return req, false
	}

	// 检查密码是否合法
	if isValid := IsValidPassword(req.Password); !isValid {
		logger.Infoln("Invalid password:", req.Password, "for username:", req.Username)
		return req, false
	}

	return req, true
}

func RegisterParseAndValidateParams(c *gin.Context) (facade.RegisterRequest, bool) {
	req := facade.RegisterRequest{
		Username: c.Query("username"),
		Password: c.Query("password"),
	}

	// 检查用户名是否合法
	if isValid := IsValidUsername(req.Username); !isValid {
		logger.Infoln("Invalid username:", req.Username)
		return req, false
	}

	// 检查密码是否合法
	if isValid := IsValidPassword(req.Password); !isValid {
		logger.Infoln("Invalid password:", req.Password, "for username:", req.Username)
		return req, false
	}

	return req, true
}

func GetUserInfoParseAndValidateParams(c *gin.Context) (int64, bool) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64) // 字符串转换为int64
	if err != nil {
		logger.Errorln("Error parsing user ID:", err)
		return 0, false
	}

	return userId, true
}

func IsValidUsername(username string) bool {
	// 用户名长度限制为3-12个字符
	const minUsernameLength = 3
	const maxUsernameLength = 32
	length := len(username)

	// 检查长度是否在范围内
	if length < minUsernameLength || length > maxUsernameLength {
		return false
	}

	// 检查用户名是否只包含字母和数字
	for _, ch := range username {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

func IsValidPassword(password string) bool {
	// 密码长度限制为3-12个字符
	const minPasswordLength = 5
	const maxPasswordLength = 32
	length := len(password)

	if length < minPasswordLength || length > maxPasswordLength {
		return false
	}

	// 密码只包括 ASCII 字母、数字和标点符号
	for _, ch := range password {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && !unicode.IsPunct(ch) {
			return false
		}
	}

	return true
}
