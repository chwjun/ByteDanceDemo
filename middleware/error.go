// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"net/http"
)

func ErrorMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			slog.Errorf("%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err,
			})
		}
	}()
	c.Next()
}
