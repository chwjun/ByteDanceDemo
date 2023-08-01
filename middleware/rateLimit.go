// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimitMiddleware
//
//	@Description: 令牌桶算法限流器
//	@param fillInterval 令牌生成间隔
//	@param cap 令牌桶容量
//	@param quantum 令牌生成量 个/次
//	@return gin.HandlerFunc
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "rate limit...",
			})
			slog.Warn("请求过于频繁")
			c.Abort()
			return
		}
		c.Next()
	}
}
