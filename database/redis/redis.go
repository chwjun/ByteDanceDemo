// Package redis @Author: youngalone [2023/8/14]
package redis

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	RateLimitClient *redis.Client
)

func Init() {
	addr := viper.GetString("settings.redis.addr")
	RateLimitClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: viper.GetString("settings.redis.password"),
		DB:       0,
	})
	if RateLimitClient == nil {
		zap.L().Error("redis初始化失败")
	} else {
		zap.L().Debug("redis初始化成功", zap.String("Addr", addr))
	}
}
