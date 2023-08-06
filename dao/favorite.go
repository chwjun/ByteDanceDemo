package dao

import (
	"errors"
	"fmt"

	"github.com/gookit/slog"

	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

func Like(userID uint, videoID uint) error {
	like := model.Like{}
	if err := DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			like = model.Like{
				UserID:  userID,
				VideoID: videoID,
				Liked:   1,
			}
			return DB.Create(&like).Error
		} else {
			return err
		}
	}

	if like.Liked == 1 {
		slog.Fatal("user has already liked this video")
		return fmt.Errorf("user has already liked this video")
	}

	like.Liked = 1
	return DB.Save(&like).Error
}

func Unlike(userID uint, videoID uint) error {
	like := model.Like{}
	if err := DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err != nil {
		return err
	}

	like.Liked = 0
	return DB.Save(&like).Error
}
