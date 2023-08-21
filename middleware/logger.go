// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func LoggerMiddleware(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	endTime := time.Now()
	durationTime := endTime.Sub(startTime)
	zap.L().Info("",
		zap.String("Method", c.Request.Method),
		zap.Int("Status", c.Writer.Status()),
		zap.Duration("durationTime", durationTime),
		zap.String("IP", c.ClientIP()),
		zap.String("URI", c.Request.RequestURI),
	)
}
