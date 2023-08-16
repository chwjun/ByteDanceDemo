package util

import (
	"fmt"

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

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *RedisClient) IncrementLikes(videoID uint) error {
	key := fmt.Sprintf("likes:%d", videoID)
	return r.client.Incr(r.ctx, key).Err()
}

func (r *RedisClient) DecrementLikes(videoID uint) error {
	key := fmt.Sprintf("likes:%d", videoID)
	return r.client.Decr(r.ctx, key).Err()
}

func (r *RedisClient) GetLikes(videoID uint) (int64, error) {
	key := fmt.Sprintf("likes:%d", videoID)
	return r.client.Get(r.ctx, key).Int64()
}
func (r *RedisClient) SyncLikesToDatabase(syncFunc func(videoID uint, likes int64) error) error {
	videoIDs := getAllVideoIDs() // 请确保此函数已定义或更改为适当的逻辑
	for _, videoID := range videoIDs {
		likes, err := r.GetLikes(videoID)
		if err != nil {
			return err
		}
		if err := syncFunc(videoID, likes); err != nil {
			return err
		}
	}
	return nil
}
