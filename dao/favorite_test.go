package dao

import (
	"testing"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/stretchr/testify/assert"
)

func TestLikeVideo(t *testing.T) {
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
		err := LikeVideo(tc.like.UserID, tc.like.VideoID)
		if tc.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)

			// 验证数据库中的值
			like := model.Like{}
			err = DB.Where(Like.UserID.Eq(tc.like.UserID), Like.VideoID.Eq(tc.like.VideoID)).First(&like).Error
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
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  2,
				VideoID: 2,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  1,
				VideoID: 9999,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  9999,
				VideoID: 1,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  20,
				VideoID: 20,
			},
			false, // 预期没有错误，因为用户已经点赞了这个视频
		},
		{
			model.Like{
				UserID:  20,
				VideoID: 22,
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

func TestGetLikedVideoIDs(t *testing.T) {
	// Prepare test data
	userID := uint(10)
	videoID1 := uint(101)
	videoID2 := uint(102)

	DB.Create(&model.Like{UserID: userID, VideoID: videoID1, Liked: 1})
	DB.Create(&model.Like{UserID: userID, VideoID: videoID2, Liked: 1})

	// Call the function under test
	videoIDs, err := GetLikedVideoIDs(userID)
	assert.NoError(t, err)

	// Check the results
	expectedVideoIDs := []uint{videoID1, videoID2}
	assert.ElementsMatch(t, expectedVideoIDs, videoIDs)
}
