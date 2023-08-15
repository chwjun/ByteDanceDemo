package main

import (
	"fmt"

	"github.com/gookit/slog"

	"github.com/RaymondCode/simple-demo/dao"
)

//	func LikeVideo(userID uint, videoID uint) error {
//		like := model.Like{}
//		if err := DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err != nil {
//			if errors.Is(err, gorm.ErrRecordNotFound) {
//				like = model.Like{
//					UserID:  userID,
//					VideoID: videoID,
//					Liked:   1,
//				}
//				return DB.Create(&like).Error
//			} else {
//				return err
//			}
//		}
//
//		if like.Liked == 1 {
//			slog.Error("user has already liked this video")
//			return fmt.Errorf("user has already liked this video")
//		}
//
//		like.Liked = 1
//		return DB.Save(&like).Error
//	}
func LikeVideo(userID uint, videoID uint) error {

	// 使用特定的查询构造方式
	// SELECT * FROM likes WHERE user_id = userID AND video_id = videoID;
	dao.SetDefault(dao.DB)
	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
	if err != nil {
		return err
	}
	println(first.Liked)
	// 假设 first 是一个 *model.Like 类型

	if first.Liked == 1 {
		slog.Error("user has already liked this video")
		return fmt.Errorf("user has already liked this video")
	}

	first.Liked = 1
	return dao.DB.Save(&first).Error
}

func main() {
	LikeVideo(1, 1)
}
