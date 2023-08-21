// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func ErrorMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("捕获错误", zap.Error(err.(error)))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err,
			})
		}
	}()
	c.Next()
}
