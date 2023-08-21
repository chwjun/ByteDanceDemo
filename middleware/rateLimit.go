// Package middleware @Author: youngalone [2023/8/1]
package middleware

import (
	redis2 "bytedancedemo/database/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func RateLimitMiddleware(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_" + ip + "_" + strconv.Itoa(int(time.Now().Unix()))
		result, err := redis2.RateLimitClient.Incr(key).Result()
		redis2.RateLimitClient.Expire(key, time.Minute)
		if err != nil {
			zap.L().Error("redis连接失败", zap.Error(err))
			return
		}
		if result > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": http.StatusTooManyRequests,
				"msg":  "请求过于频繁",
			})
			zap.L().Warn("请求过于频繁")
			c.Abort()
			return
		}
		c.Next()
	}
}
