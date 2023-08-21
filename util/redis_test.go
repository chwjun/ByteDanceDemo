package util

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkGetTotalLikeCounts(b *testing.B) {

	videoIDs := []uint{1, 2, 3}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GlobalRedisClient.GetTotalLikeCounts(videoIDs)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func TestGetVideosLikes(t *testing.T) {

	// 定义测试数据
	videoIDs := []uint{1, 2, 3, 4}
	expectedLikes := map[uint]int64{
		1: 100,
		2: 200,
		3: 300,
	}

	// 将测试数据插入Redis
	for videoID, likes := range expectedLikes {
		videoKey := fmt.Sprintf("video:%d", videoID)
		videoLikesField := "totalVideoLikes"
		GlobalRedisClient.client.HSet(GlobalRedisClient.ctx, videoKey, videoLikesField, strconv.FormatInt(likes, 10))
	}

	// 调用函数
	result, err := GetVideosLikes(videoIDs)
	if err != nil {
		t.Fatalf("GetVideosLikes failed: %v", err)
	}
	// 打印返回的结果
	t.Logf("Result: %+v", result)

	// 验证结果
	for videoID, likes := range expectedLikes {
		if gotLikes, ok := result[videoID]; !ok || gotLikes != likes {
			t.Errorf("Expected %d likes for videoID %d, got %d", likes, videoID, gotLikes)
		}
	}

}
