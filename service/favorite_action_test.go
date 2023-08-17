package service

import (
	"sync"
	"testing"
	_ "time"

	"github.com/stretchr/testify/assert"
)

func TestFavoriteActionIntegration(t *testing.T) {

	// 创建服务实例
	favoriteService := NewFavoriteService("43.140.203.85:6379", "", 0)

	// 调用方法
	userId := uint(2)
	videoID := int64(10)
	actionType := int32(1)
	response, err := favoriteService.FavoriteAction(userId, videoID, actionType)

	// 验证结果
	assert.NoError(t, err, "FavoriteAction should not return an error")
	assert.Equal(t, SuccessCode, response.StatusCode, "Unexpected status code")
	// 添加其他验证
}

func BenchmarkFavoriteAction(b *testing.B) {
	// 创建服务实例
	favoriteService := NewFavoriteService("43.140.203.85:6379", "", 0)

	// 设置参数
	userId := uint(2)
	videoID := int64(10)
	actionType := int32(1)

	var wg sync.WaitGroup

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			favoriteService.FavoriteAction(userId, videoID, actionType)
			wg.Done()
		}()
	}

	wg.Wait() // 等待所有的FavoriteAction函数调用都完成
}
