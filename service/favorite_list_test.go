package service

import (
	"testing"
)

func BenchmarkFavoriteList(b *testing.B) {
	// 创建服务实例
	service := &FavoriteServiceImpl{
		// 您可以在此处设置任何所需的模拟或实际依赖项
	}

	// 设置测试数据
	userID := int64(1)

	// 运行基准测试
	b.ResetTimer() // 重置计时器以排除设置过程的时间
	for i := 0; i < b.N; i++ {
		_, _ = service.FavoriteList(userID)
	}
}
