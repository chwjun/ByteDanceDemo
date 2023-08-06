package dao

import (
	"testing"
	"time"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/stretchr/testify/assert"
)

func TestLike(t *testing.T) {
	// 定义你的测试用例，包括一个Like对象以及预期的结果
	testCases := []struct {
		like    model.Like
		wantErr bool
	}{
		{
			model.Like{
				BaseModel: model.BaseModel{
					ID:        10,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				UserID:  9999, // 不存在的 UserID
				VideoID: 1,    // 存在的 VideoID
				Liked:   1,
			},
			true, // 预期有错误，因为用户ID不存在
		},
		// 测试用例2
		struct {
			like    model.Like
			wantErr bool
		}{
			model.Like{
				BaseModel: model.BaseModel{
					ID:        11,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				UserID:  1,    // 存在的 UserID
				VideoID: 9999, // 不存在的 VideoID
				Liked:   1,
			},
			true, // 预期有错误，因为视频ID不存在
		},
		// 测试用例3
		struct {
			like    model.Like
			wantErr bool
		}{
			model.Like{
				BaseModel: model.BaseModel{
					ID:        3,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				UserID:  1, // 存在的 UserID
				VideoID: 1, // 存在的 VideoID
				Liked:   1,
			},
			true, // 预期有错误，因为用户已经点赞了这个视频
		},
		// 测试用例4
		struct {
			like    model.Like
			wantErr bool
		}{
			model.Like{
				BaseModel: model.BaseModel{
					ID:        4,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				UserID:  2, // 存在的 UserID
				VideoID: 2, // 存在的 VideoID
				Liked:   0,
			},
			false, // 预期没有错误，因为用户还没有点赞这个视频
		},
	}

	for _, tc := range testCases {
		err := Like(tc.like.UserID, tc.like.VideoID)
		if tc.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)

			// 验证数据库中的值
			like := model.Like{}
			err = DB.Where("user_id = ? AND video_id = ?", tc.like.UserID, tc.like.VideoID).First(&like).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, like.Liked)
		}
	}
}
