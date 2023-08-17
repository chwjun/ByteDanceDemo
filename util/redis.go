package util

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// 检查连接状态
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	} else {
		log.Printf("Successfully connected to Redis at %s", addr)
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *RedisClient) LikeVideo(userID uint, videoID uint) error {
	key := fmt.Sprintf("likes:%d", videoID)
	field := fmt.Sprintf("%d", userID)
	return r.client.HSet(r.ctx, key, field, 1).Err()
}

func (r *RedisClient) UnlikeVideo(userID uint, videoID uint) error {
	key := fmt.Sprintf("likes:%d", videoID)
	field := fmt.Sprintf("%d", userID)
	return r.client.HSet(r.ctx, key, field, 0).Err()
}

func (r *RedisClient) GetLikes(videoID uint) (map[string]string, error) {
	key := fmt.Sprintf("likes:%d", videoID)
	return r.client.HGetAll(r.ctx, key).Result()
}
