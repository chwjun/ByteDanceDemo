package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkFavoriteAction(b *testing.B) {
	service := &FavoriteServiceImpl{}
	userId := int64(123)
	videoID := int64(456)
	actionType := int32(1)

	// 重置计时器以忽略基准测试的初始化部分
	b.ResetTimer()

	// 使用 b.RunParallel 运行并发基准测试
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := service.FavoriteAction(userId, videoID, actionType)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
func TestLikeVideoAlreadyLiked(t *testing.T) {
	// 假设的用户和视频ID
	userID := int64(1)
	videoID := int64(181)

	// 先确保点赞
	err := likeVideo(userID, videoID)
	assert.Nil(t, err) // 断言没有错误

	// 再次尝试点赞同一视频，应该返回 "user has already liked this video" 错误
	err = likeVideo(userID, videoID)
	assert.NotNil(t, err)                                             // 断言有错误
	assert.Equal(t, "user has already liked this video", err.Error()) // 断言错误消息
}
