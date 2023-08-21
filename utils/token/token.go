// Package token @Author: youngalone [2023/8/4]
package token

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	UserID   int64  `json:"user_id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}

func GenerateToken(secretKey []byte, claims Claims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(secretKey []byte, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token令牌非法")
	}
	return token.Claims.(*Claims), err
}
