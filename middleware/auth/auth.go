package auth

import (
	"log"
	"net/http"
	"os"

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
	}
	if len(signaturedString) == 0 {
		c.Abort()
		c.JSON(
			http.StatusUnauthorized,
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
		log.Println("Token Right:", claims.Name, claims.Id)
		c.Set("username", claims.Name)
		c.Set("id", claims.Id)
		c.Next()
	default:
		log.Println("Token Error.")
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			Response{
				StatusCode: -1,
				StatusMsg:  "Token Error.",
			})
	}
}

// ----------------------------------private-------------------------------------

var signatureVarName = "TIKTOKLITESIGNINGKEY"

func parseKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv(signatureVarName)), nil
}
