package dao

import (
	"errors"
	"fmt"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/gookit/slog"
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
		slog.Error("user has already liked this video")
		return fmt.Errorf("user has already liked this video")
	}

	like.Liked = 1
	return DB.Save(&like).Error
}

func Unlike(userID uint, videoID uint) error {
	like := model.Like{}
	if err := DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("No like found for this user and video")
		}
		return err
	}

	if like.Liked == 0 {
		slog.Error("User has already unliked this video")
		return fmt.Errorf("User has already unliked this video")
	}

	like.Liked = 0
	return DB.Save(&like).Error
}

// 获取用户点赞的所有视频ID
func GetLikedVideoIDs(userID uint) ([]uint, error) {
	var videoIDs []uint
	err := DB.Model(&model.Like{}).Where("user_id = ? AND liked = 1", userID).Pluck("video_id", &videoIDs).Error
	if err != nil {
		return nil, err
	}
	return videoIDs, nil
}

// 通过视频ID获取视频详情
func GetVideoDetailsByID(videoID uint) (uint, string, string, string, error) {
	video := model.Video{}
	if err := DB.Select("author_id", "title", "play_url", "cover_url").Where("id = ?", videoID).First(&video).Error; err != nil {
		return 0, "", "", "", err
	}
	return video.AuthorID, video.Title, video.PlayURL, video.CoverURL, nil
}

// 通过视频ID获取评论总数
func GetCommentCount(videoID uint) (int64, error) {
	var count int64
	if err := DB.Model(&model.Comment{}).Where("video_id = ?", videoID).Where("action_type = ?", 1).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 通过视频ID获取点赞总数
func GetLikeCount(videoID uint) (int64, error) {
	var count int64
	if err := DB.Model(&model.Like{}).Where("video_id = ?", videoID).Where("liked = ?", 1).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 获取某用户关注总数
func GetUserFollowCount(userID uint) (int64, error) {
	var followCount int64
	if err := DB.Model(&model.Relation{}).Where("user_id = ?", userID).Where("followed = ?", 1).Count(&followCount).Error; err != nil {
		return 0, err
	}
	return followCount, nil
}

// 获取某用户粉丝总数
func GetUserFollowerCount(userID uint) (int64, error) {
	var followerCount int64
	if err := DB.Model(&model.Relation{}).Where("following_id = ?", userID).Where("followed = ?", 1).Count(&followerCount).Error; err != nil {
		return 0, err
	}
	return followerCount, nil
}

// 检查一个用户是否关注了另一个用户
func IsUserFollowingAnotherUser(userID, followingID uint) (bool, error) {
	var relation model.Relation
	if err := DB.Where("user_id = ?", userID).Where("following_id = ?", followingID).Where("followed = ?", 1).First(&relation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // 没有找到关系记录，表示用户没有关注另一个用户
		}
		return false, err // 其他错误
	}

	return true, nil // 找到关系记录，表示用户关注了另一个用户
}

// 获取用户的获赞总数
func GetUserTotalReceivedLikes(userID uint) (int64, error) {
	var totalLikes int64
	if err := DB.Table("likes").Joins("JOIN videos on videos.id = likes.video_id").Where("videos.author_id = ?", userID).Count(&totalLikes).Error; err != nil {
		return 0, err
	}
	return totalLikes, nil
}

// 获取用户的点赞总数（用户给出的赞）
func GetUserTotalGivenLikes(userID uint) (int64, error) {
	var totalLikes int64
	if err := DB.Model(&model.Like{}).Where("user_id = ?", userID).Where("liked = ?", 1).Count(&totalLikes).Error; err != nil {
		return 0, err
	}
	return totalLikes, nil
}

// 获取用户的作品数量
func GetUserWorkCount(userID uint) (int64, error) {
	var workCount int64
	if err := DB.Model(&model.Video{}).Where("author_id = ?", userID).Count(&workCount).Error; err != nil {
		return 0, err
	}
	return workCount, nil
}
