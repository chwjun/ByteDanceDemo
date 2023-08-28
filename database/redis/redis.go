// Package redis @Author: youngalone [2023/8/14]
package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var (
	Addr       string
	Password   string
	ExpireTime time.Duration
	Ctx        context.Context
	// RdbTest 测试链接
	RdbTest *redis.Client
	// RateLimitClient 限流中间件链接
	RateLimitClient *redis.Client
	// UserFollowings 根据用户id找到他关注的人
	UserFollowings *redis.Client
	// UserFollowers 根据用户id找到他的粉丝
	UserFollowers *redis.Client
	// UserFriends 根据用户id找到他的好友
	UserFriends *redis.Client
	// RdbVCid 存储video与comment的关系
	RdbVCid *redis.Client
	// RdbCVid 根据commentId找videoId
	RdbCVid *redis.Client
	// RdbCIdComment 根据commentId 找comment
	RdbCIdComment *redis.Client
)

func Init() {
	Addr = fmt.Sprintf("%s:%s",
		viper.GetString("settings.redis.host"),
		viper.GetString("settings.redis.port"),
	)
	Password = viper.GetString("settings.redis.password")
	ExpireTime = viper.GetDuration("settings.redis.expirationTime") * time.Second

	Ctx = context.Background()

	RdbTest = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       0,
	})
	RateLimitClient = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       1,
	})
	RdbVCid = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       2,
	})
	RdbCVid = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       3,
	})
	RdbCIdComment = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       4,
	})
	UserFollowings = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       5,
	})
	UserFollowers = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       6,
	})
	UserFriends = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       7,
	})
	_, err := RdbTest.Ping(Ctx).Result()
	if err != nil {
		zap.L().Error("redis初始化失败")
	} else {
		zap.L().Debug("redis初始化成功", zap.String("Addr", Addr))
	}
}
