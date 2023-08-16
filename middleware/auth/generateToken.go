package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	Name       string `json:"name"` // 用户名
	Id         int64  `json:"id"`   // id
	GenerateAt int64  // 生成时的UNIX时间戳
	jwt.RegisteredClaims
}

func GenerateToken(name string, id int64) (string, error) {
	log.Printf("Generating token for name:%v id:%v\n", name, id)

	var (
		SIGNINGKEY = []byte(os.Getenv(signatureVarName)) // 密钥
		claims     = &customClaims{
			name,
			id,
			time.Now().Unix(),
			jwt.RegisteredClaims{},
		}
	)

	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signaturedString, err := userToken.SignedString(SIGNINGKEY); err != nil {
		log.Println("Failed." + err.Error())
		return "", err
	} else {
		log.Println("Succeeded.")
		return signaturedString, nil
	}
}
