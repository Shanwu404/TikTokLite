package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func Auth(c *gin.Context) {
	signaturedString := c.Query("token")
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
		return
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
	if !osSsignaturedStringValid() {
		err := errors.New(fmt.Sprintf("The %v is blank!!!", signatureVarName))
		setTempSignature()
		log.Println(err)
		return []byte{}, err
	}
	return []byte(os.Getenv(signatureVarName)), nil
}

func osSsignaturedStringValid() bool {
	return os.Getenv(signatureVarName) != ""
}

func setTempSignature() error {
	err := os.Setenv(signatureVarName, fmt.Sprint((time.Now().Nanosecond())))
	if err != nil {
		log.Println("Error executing command:", err)
		return err
	} else {
		log.Printf("Set %v as tempKey.\n", signatureVarName)
	}
	return nil
}
