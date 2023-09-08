package auth

import (
	"net/http"
	"os"

	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// TODO:检查token里含有特殊字符如% 是否需要特殊处理或者统一使用合适的编解码方法
func Auth(c *gin.Context) {
	var signaturedString string
	switch c.Request.Method {
	case "GET":
		signaturedString = c.Query("token")
	case "POST":
		signaturedString = c.PostForm("token")
		if len(signaturedString) == 0 {
			signaturedString = c.Query("token")
		}
	}
	if len(signaturedString) == 0 && c.Request.URL.Path == "/douyin/feed/" {
		logger.Infoln("Feeding without token.")
		return
	}
	if len(signaturedString) == 0 {
		c.Abort()
		c.JSON(
			// 这里用401客户端会显示网络错误，影响未登录用户体验
			// http.StatusUnauthorized,
			http.StatusOK,
			Response{
				StatusCode: -1,
				StatusMsg:  "Please verify the login status.",
			})
		return
	}

	decodedToken, err := jwt.ParseWithClaims(
		signaturedString,
		&customClaims{},
		parseKeyFunc,
	)

	claims, ok := decodedToken.Claims.(*customClaims)
	switch {
	case err == nil && ok && decodedToken.Valid:
		logger.Infoln("Valid Token:", claims.Name, claims.Id)
		c.Set("username", claims.Name)
		c.Set("id", claims.Id)
		return
	default:
		logger.Infoln("Invalid token.")
		c.AbortWithStatusJSON(
			// 这里用401客户端会显示网络错误，影响未登录用户体验
			// http.StatusUnauthorized,
			http.StatusOK,
			Response{
				StatusCode: -1,
				StatusMsg:  "Invalid token.",
			})
	}
}

// ----------------------------------private-------------------------------------

var signatureVarName = "TIKTOKLITESIGNINGKEY"

func parseKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv(signatureVarName)), nil
}
