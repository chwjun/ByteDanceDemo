// Package middleware @Author: youngalone [2023/8/6]
package middleware

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func JWTMiddleware(c *gin.Context) {
	path := c.Request.URL.Path
	if path == "/douyin/user/register/" || path == "/douyin/user/login/" {
		c.Set("role", "tourist")
		c.Next()
		return
	}
	tokenString := c.Query("token")
	if tokenString == "" {
		tokenString = c.Request.PostFormValue("token")
	}
	if tokenString == "" {
		c.Set("role", "tourist")
		c.Next()
	} else {
		claims, err := token.ParseToken([]byte(viper.GetString("settings.jwt.secretKey")), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, fmt.Sprintf("token非法 %v", err))
			c.Abort()
			return
		}
		if err = claims.Valid(); err != nil {
			c.JSON(http.StatusUnauthorized, fmt.Sprintf("token非法 %v", err))
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.UserName)
		c.Set("role", claims.Role)
		c.Next()
	}
}
