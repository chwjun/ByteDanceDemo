package dao

import (
	"testing"

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
				UserID:  9999,
				VideoID: 1,
			},
			false,
		},
		{
			model.Like{
				UserID:  1,
				VideoID: 9999,
			},
			true,
		},
		{
			model.Like{
				UserID:  1,
				VideoID: 1,
			},
			true,
		},
		{
			model.Like{
				UserID:  2,
				VideoID: 2,
			},
			false,
		}, {
			model.Like{
				UserID:  17,
				VideoID: 17,
			},
			true,
		}, {
			model.Like{
				UserID:  20,
				VideoID: 22,
			},
			true,
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

func TestUnlike(t *testing.T) {
	// 定义你的测试用例，包括一个Like对象以及预期的结果
	testCases := []struct {
		like    model.Like
		wantErr bool
	}{
		{
			model.Like{
				UserID:  1,
				VideoID: 1,
				Liked:   1,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  2,
				VideoID: 2,
				Liked:   1,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  1,
				VideoID: 9999,
				Liked:   1,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  9999,
				VideoID: 1,
				Liked:   1,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  20,
				VideoID: 20,
				Liked:   1,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  20,
				VideoID: 22,
				Liked:   1,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
	}

	for _, tc := range testCases {
		err := Unlike(tc.like.UserID, tc.like.VideoID)
		if tc.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)

			// 验证数据库中的值
			like := model.Like{}
			err = DB.Where("user_id = ? AND video_id = ?", tc.like.UserID, tc.like.VideoID).First(&like).Error
			assert.NoError(t, err)
			assert.Equal(t, 0, like.Liked)
		}
	}
}
