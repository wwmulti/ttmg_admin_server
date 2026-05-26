package utils

import (
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int64, username string) (string, error) {
	claims := CustomClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    beego.AppConfig.DefaultString("appname", "api"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(beego.AppConfig.DefaultString("JwtSecretKey", "default")))
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.DefaultString("JwtSecretKey", "default")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
