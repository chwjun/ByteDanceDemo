package utils

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"

	"github.com/gookit/slog"

	"golang.org/x/net/context"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var GlobalRedisClient *RedisClient

func init() {
	// 你的 Redis 配置
	addr := "43.140.203.85:6388"
	password := "sample_douyin"
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
func UpdateLikeCounts(userID int64, videoID int64, like bool) error {
	// 定义哈希表的键
	userKey := fmt.Sprintf("user:%d", userID)
	videoKey := fmt.Sprintf("video:%d", videoID)
	slog.Debug("updateLikeCounts")
	// 定义操作数（增加或减少）
	var operation int64
	if like {
		operation = 1
	} else {
		operation = -1
	}

	// 获取 Redis 上下文
	ctx := GlobalRedisClient.ctx
	pipe := GlobalRedisClient.client.Pipeline()

	// 更新用户的点赞总数
	userLikesField := "totalLikes"
	userUpdateCmd := pipe.HIncrBy(ctx, userKey, userLikesField, operation)

	// 更新视频的获赞总数
	videoLikesField := "totalVideoLikes"
	videoUpdateCmd := pipe.HIncrBy(ctx, videoKey, videoLikesField, operation)

	// 执行管道中的所有命令
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute pipeline: %v", err)
	}

	// 检查用户更新的错误
	if err := userUpdateCmd.Err(); err != nil {
		return fmt.Errorf("failed to update user likes: %v", err)
	}

	// 检查视频更新的错误
	if err := videoUpdateCmd.Err(); err != nil {
		return fmt.Errorf("failed to update video likes: %v", err)
	}

	return nil
}
func GetVideosLikes(videoIDs []int64) (map[int64]int64, error) {
	ctx := GlobalRedisClient.ctx
	pipe := GlobalRedisClient.client.Pipeline()

	futures := make(map[int64]*redis.StringCmd)
	for _, videoID := range videoIDs {
		videoKey := fmt.Sprintf("video:%d", videoID)
		videoLikesField := "totalVideoLikes"
		futures[videoID] = pipe.HGet(ctx, videoKey, videoLikesField)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	result := make(map[int64]int64)
	for videoID, future := range futures {
		err := future.Err()
		if err == redis.Nil {
			result[videoID] = 0
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get video likes for video %d: %v", videoID, err)
		}
		likesStr, _ := future.Result()
		likes, err := strconv.ParseInt(likesStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse likes for video %d: %v", videoID, err)
		}
		result[videoID] = likes
	}

	return result, nil
}

func GetUserFavorites(userIDs []int64) (map[int64]int64, error) {
	ctx := GlobalRedisClient.ctx
	pipe := GlobalRedisClient.client.Pipeline()

	// 创建一个存储未来结果的映射
	futures := make(map[int64]*redis.StringCmd)
	for _, userID := range userIDs {
		userKey := fmt.Sprintf("user:%d", userID)
		userLikesField := "totalLikes"
		futures[userID] = pipe.HGet(ctx, userKey, userLikesField)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	// 提取结果
	result := make(map[int64]int64)
	for userID, future := range futures {
		likesStr, err := future.Result()
		if err != nil {
			if err == redis.Nil {
				// 如果用户不存在，则设置键并将值设为0
				userKey := fmt.Sprintf("user:%d", userID)
				userLikesField := "totalLikes"
				if err := GlobalRedisClient.client.HSet(ctx, userKey, userLikesField, "0").Err(); err != nil {
					return nil, fmt.Errorf("failed to set user favorites for user %d: %v", userID, err)
				}
				continue
			}
			return nil, fmt.Errorf("failed to get user favorites for user %d: %v", userID, err)
		}
		likes, err := strconv.ParseInt(likesStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse favorites for user %d: %v", userID, err)
		}
		result[userID] = likes
	}

	return result, nil
}
func GetTotalVideosLikes(videoIDs []int64) (int64, error) {
	ctx := GlobalRedisClient.ctx
	pipe := GlobalRedisClient.client.Pipeline()

	// 创建一个存储未来结果的切片
	futures := make([]*redis.StringCmd, len(videoIDs))

	// 遍历每个视频 ID，将 HGet 命令添加到管道
	for i, videoID := range videoIDs {
		videoKey := fmt.Sprintf("video:%d", videoID)
		videoLikesField := "totalVideoLikes"
		futures[i] = pipe.HGet(ctx, videoKey, videoLikesField)
	}

	// 一次性执行所有命令
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	// 总和变量，用于存储所有视频的喜欢总数
	totalLikes := int64(0)

	// 提取每个未来结果，并累加到总和
	for _, future := range futures {
		likesStr, err := future.Result()
		if err != nil {
			if err == redis.Nil {
				// 如果视频不存在，则跳过
				continue
			}
			return 0, fmt.Errorf("failed to get video likes: %v", err)
		}
		likes, err := strconv.ParseInt(likesStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse likes: %v", err)
		}
		totalLikes += likes
	}

	return totalLikes, nil
}

//func (r *RedisClient) LikeVideo(userID int64, videoID int64) error {
//	key := fmt.Sprintf("likes:%d", videoID)
//	field := fmt.Sprintf("%d", userID)
//	totalLikesField := "totalLikes"
//
//	// 获取当前点赞总数
//	currentTotalLikes, err := r.client.HGet(r.ctx, key, totalLikesField).Int64()
//	if err != nil {
//		if err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
//			return err
//		}
//		currentTotalLikes = 0 // 如果键不存在，设置当前点赞总数为0
//	}
//
//	// 事务开始
//	pipe := r.client.TxPipeline()
//	pipe.HSet(r.ctx, key, field, 1)                             // 存储用户点赞
//	pipe.HSet(r.ctx, key, totalLikesField, currentTotalLikes+1) // 更新点赞总数
//
//	// 执行事务
//	_, err = pipe.Exec(r.ctx)
//	return err
//}

//func (r *RedisClient) UnlikeVideo(userID int64, videoID int64) error {
//	key := fmt.Sprintf("likes:%d", videoID)
//	field := fmt.Sprintf("%d", userID)
//	totalLikesField := "totalLikes"
//
//	// 检查用户是否已经点赞
//	userLike, err := r.client.HGet(r.ctx, key, field).Int64()
//	if err != nil {
//		if err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
//			return err
//		}
//		// 如果键不存在或用户未点赞，则无需执行任何操作
//		return nil
//	}
//
//	// 如果用户未点赞，则无需执行任何操作
//	if userLike == 0 {
//		return nil
//	}
//
//	// 获取当前点赞总数
//	currentTotalLikes, err := r.client.HGet(r.ctx, key, totalLikesField).Int64()
//	if err != nil {
//		if err != redis.Nil { // 如果错误不是由于键不存在造成的，则返回错误
//			return err
//		}
//		currentTotalLikes = 0 // 如果键不存在，设置当前点赞总数为0
//	}
//
//	// 事务开始
//	pipe := r.client.TxPipeline()
//	pipe.HSet(r.ctx, key, field, 0)                             // 存储用户取消点赞
//	pipe.HSet(r.ctx, key, totalLikesField, currentTotalLikes-1) // 更新点赞总数
//
//	// 执行事务
//	_, err = pipe.Exec(r.ctx)
//	return err
//}

func (r *RedisClient) GetVideoLikeCounts(videoIDs []int64) (map[int64]int64, error) {
	resultsMap := make(map[int64]int64)

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
func (r *RedisClient) GetTotalLikeCounts(videoIDs []int64) (int64, error) {
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

//func SyncLikesToDatabase() error {
//	videoIDs, err := getAllVideoIDs() // 获取所有视频ID
//	if err != nil {
//		return err
//	}
//
//	for _, videoID := range videoIDs {
//		likes, err := GlobalRedisClient.GetLikes(videoID)
//		if err != nil {
//			return err
//		}
//
//		for userIDString, likedString := range likes {
//			userID, err := strconv.Atoi(userIDString)
//			if err != nil {
//				return err // 或者可以记录错误并继续
//			}
//			liked, err := strconv.Atoi(likedString)
//			if err != nil {
//				return err // 或者可以记录错误并继续
//			}
//
//			// 开始一个新的事务
//			tx := dao.DB.Begin()
//			if tx.Error != nil {
//				return tx.Error
//			}
//
//			first, err := dao.Like.Where(dao.Like.UserID.Eq(int64(userID)), dao.Like.VideoID.Eq(int64(videoID))).First()
//			if err != nil {
//				if errors.Is(err, gorm.ErrRecordNotFound) {
//					// 如果记录未找到，则创建一个新的喜欢记录
//					like := model.Like{
//						UserID:  int64(userID),
//						VideoID: int64(videoID),
//						Liked:   int(int64(liked)),
//					}
//					if err := tx.Create(&like).Error; err != nil {
//						tx.Rollback() // 回滚事务
//						return err
//					}
//				} else {
//					// 如果发生其他错误，则回滚事务并返回该错误
//					tx.Rollback()
//					return err
//				}
//			} else {
//				// 如果记录存在，则更新点赞状态
//				first.Liked = int(int64(liked))
//				if err := tx.Save(&first).Error; err != nil {
//					tx.Rollback() // 回滚事务
//					return err
//				}
//			}
//
//			// 提交事务
//			if err := tx.Commit().Error; err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//func getAllVideoIDs() ([]int64, error) {
//	videos := dao.Video // 假设您有一个名为dao的包，其中有一个Video的DAO对象
//
//	var videoIDs []int64
//	ids, err := videos.Select(videos.ID).Find()
//	if err != nil {
//		return nil, err
//	}
//
//	for _, id := range ids {
//		videoIDs = append(videoIDs, id.ID) // 假设ID是Video中的一个字段
//	}
//
//	return videoIDs, nil
//}
//
//func StartSyncTask(db *gorm.DB, syncInterval time.Duration) {
//	ticker := time.NewTicker(syncInterval)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ticker.C:
//			if err := SyncLikesToDatabase(); err != nil {
//				log.Printf("Failed to sync likes to database: %v", err)
//			}
//		}
//	}
//}
