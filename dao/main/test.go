package main

import (
	"fmt"

	"github.com/RaymondCode/simple-demo/dao"
)

type LikeVideoResult struct {
	VideoID uint
}

func AreVideosLikedByUser(userID uint, videoIDs []uint) (map[uint]bool, error) {
	likedVideos := make(map[uint]bool)

	var results []LikeVideoResult
	err := dao.Like.Select(dao.Like.VideoID).Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.In(videoIDs...), dao.Like.Liked.Eq(1), dao.Like.DeletedAt.IsNull()).Scan(&results)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户喜欢的视频: %v", err)
	}

	for _, result := range results {
		likedVideos[result.VideoID] = true
	}

	// Set false for videos not liked by the user
	for _, id := range videoIDs {
		if _, ok := likedVideos[id]; !ok {
			likedVideos[id] = false
		}
	}

	return likedVideos, nil
}

func main() {
	// 用户ID和视频ID列表
	userID := uint(1)
	videoIDs := []uint{1, 2, 3, 4}

	// 调用AreVideosLikedByUser函数检查用户是否喜欢了这些视频
	likedVideos, err := AreVideosLikedByUser(userID, videoIDs)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 打印结果
	fmt.Printf("User %d liked videos:\n", userID)
	for videoID, liked := range likedVideos {
		if liked {
			fmt.Printf("- Video %d (liked)\n", videoID)
		} else {
			fmt.Printf("- Video %d (not liked)\n", videoID)
		}
	}
}
