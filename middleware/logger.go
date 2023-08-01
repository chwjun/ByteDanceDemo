// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"time"
)

func LoggerMiddleware(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	endTime := time.Now()
	durationTime := endTime.Sub(startTime)
	slog.Infof("%v | %v | %v | %v | \"%v\"", c.Request.Method, c.Writer.Status(), durationTime, c.ClientIP(), c.Request.RequestURI)
}
