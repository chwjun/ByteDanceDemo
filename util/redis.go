package util

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var GlobalRedisClient *RedisClient

func init() {
	// 你的 Redis 配置
	addr := "43.140.203.85:6379"
	password := ""
	db := 0

	// 初始化全局 RedisClient
	GlobalRedisClient = NewRedisClient(addr, password, db)
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
	totalLikesField := "totalLikes"

	// 获取当前点赞总数
	currentTotalLikes, err := r.client.HGet(r.ctx, key, totalLikesField).Int64()
	if err != nil {
		if err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
			return err
		}
		currentTotalLikes = 0 // 如果键不存在，设置当前点赞总数为0
	}

	// 事务开始
	pipe := r.client.TxPipeline()
	pipe.HSet(r.ctx, key, field, 1)                             // 存储用户点赞
	pipe.HSet(r.ctx, key, totalLikesField, currentTotalLikes+1) // 更新点赞总数

	// 执行事务
	_, err = pipe.Exec(r.ctx)
	return err
}

func (r *RedisClient) UnlikeVideo(userID uint, videoID uint) error {
	key := fmt.Sprintf("likes:%d", videoID)
	field := fmt.Sprintf("%d", userID)
	totalLikesField := "totalLikes"

	// 检查用户是否已经点赞
	userLike, err := r.client.HGet(r.ctx, key, field).Int64()
	if err != nil {
		if err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
			return err
		}
		// 如果键不存在或用户未点赞，则无需执行任何操作
		return nil
	}

	// 如果用户未点赞，则无需执行任何操作
	if userLike == 0 {
		return nil
	}

	// 获取当前点赞总数
	currentTotalLikes, err := r.client.HGet(r.ctx, key, totalLikesField).Int64()
	if err != nil {
		if err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
			return err
		}
		currentTotalLikes = 0 // 如果键不存在，设置当前点赞总数为0
	}

	// 事务开始
	pipe := r.client.TxPipeline()
	pipe.HSet(r.ctx, key, field, 0)                             // 存储用户取消点赞
	pipe.HSet(r.ctx, key, totalLikesField, currentTotalLikes-1) // 更新点赞总数

	// 执行事务
	_, err = pipe.Exec(r.ctx)
	return err
}

func (r *RedisClient) GetLikes(videoID uint) (map[string]string, error) {
	key := fmt.Sprintf("likes:%d", videoID)
	return r.client.HGetAll(r.ctx, key).Result()
}
func (r *RedisClient) GetVideoLikeCounts(videoIDs []uint) (map[uint]int64, error) {
	resultsMap := make(map[uint]int64)

	for _, videoID := range videoIDs {
		likeKey := fmt.Sprintf("likes:%d", videoID)
		likeCount, err := r.client.HGet(r.ctx, likeKey, "totalLikes").Int64()
		if err != nil {
			if err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
				return nil, fmt.Errorf("无法获取视频的喜欢（like）数量: %v", err)
			}
			likeCount = 0 // 如果键不存在，设置当前点赞总数为0
		}
		resultsMap[videoID] = likeCount
	}

	return resultsMap, nil
}
func (r *RedisClient) GetTotalLikeCounts(videoIDs []uint) (int64, error) {
	// 初始化总点赞数为 0
	totalLikeCount := int64(0)

	// 遍历每个视频ID，并获取其点赞数
	for _, videoID := range videoIDs {
		likeKey := fmt.Sprintf("likes:%d", videoID)
		likeCount, err := r.client.HGet(r.ctx, likeKey, "totalLikes").Int64()
		if err != nil && err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
			return 0, fmt.Errorf("无法获取视频的喜欢（like）数量: %v", err)
		}
		totalLikeCount += likeCount // 将找到的点赞数加到总数上
	}

	return totalLikeCount, nil
}

func SyncLikesToDatabase() error {
	videoIDs, err := getAllVideoIDs() // 获取所有视频ID
	if err != nil {
		return err
	}

	for _, videoID := range videoIDs {
		likes, err := GlobalRedisClient.GetLikes(videoID)
		if err != nil {
			return err
		}

		for userIDString, likedString := range likes {
			userID, err := strconv.Atoi(userIDString)
			if err != nil {
				return err // 或者可以记录错误并继续
			}
			liked, err := strconv.Atoi(likedString)
			if err != nil {
				return err // 或者可以记录错误并继续
			}

			// 开始一个新的事务
			tx := dao.DB.Begin()
			if tx.Error != nil {
				return tx.Error
			}

			first, err := dao.Like.Where(dao.Like.UserID.Eq(uint(userID)), dao.Like.VideoID.Eq(uint(videoID))).First()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 如果记录未找到，则创建一个新的喜欢记录
					like := model.Like{
						UserID:  uint(userID),
						VideoID: uint(videoID),
						Liked:   int(uint(liked)),
					}
					if err := tx.Create(&like).Error; err != nil {
						tx.Rollback() // 回滚事务
						return err
					}
				} else {
					// 如果发生其他错误，则回滚事务并返回该错误
					tx.Rollback()
					return err
				}
			} else {
				// 如果记录存在，则更新点赞状态
				first.Liked = int(uint(liked))
				if err := tx.Save(&first).Error; err != nil {
					tx.Rollback() // 回滚事务
					return err
				}
			}

			// 提交事务
			if err := tx.Commit().Error; err != nil {
				return err
			}
		}
	}
	return nil
}
func getAllVideoIDs() ([]uint, error) {
	videos := dao.Video // 假设您有一个名为dao的包，其中有一个Video的DAO对象

	var videoIDs []uint
	ids, err := videos.Select(videos.ID).Find()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		videoIDs = append(videoIDs, id.ID) // 假设ID是Video中的一个字段
	}

	return videoIDs, nil
}

func StartSyncTask(db *gorm.DB, syncInterval time.Duration) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := SyncLikesToDatabase(); err != nil {
				log.Printf("Failed to sync likes to database: %v", err)
			}
		}
	}
}
