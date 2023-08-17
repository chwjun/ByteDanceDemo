package service

import (
	"testing"
)

func BenchmarkFavoriteAction(b *testing.B) {
	service := &FavoriteServiceImpl{}
	userId := int64(123)
	videoID := int64(456)
	actionType := int32(1)

	// 重置计时器以忽略基准测试的初始化部分
	b.ResetTimer()

	// N是基准测试的迭代次数; b.N会在测试运行期间自动调整
	for i := 0; i < b.N; i++ {
		_, err := service.FavoriteAction(userId, videoID, actionType)
		if err != nil {
			b.Fatal(err)
		}
	}
}
