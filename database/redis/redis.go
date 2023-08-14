// Package redis @Author: youngalone [2023/8/14]
package redis

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	RateLimitClient *redis.Client
)

func Init() {
	RateLimitClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("settings.redis.addr"),
		Password: viper.GetString("settings.redis.password"),
		DB:       0,
	})
}
