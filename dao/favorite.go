package dao

//
//import (
//	"fmt"
//
//	"github.com/RaymondCode/simple-demo/model"
//	"gorm.io/gorm"
//)
//
//func GetUserDetailsByID(userID uint) (string, string, string, string, error) {
//	user := model.User{}
//	if err := DB.Select("name, avatar, background_image, signature").Where("id = ?", userID).First(&user).Error; err != nil {
//		return "", "", "", "", err
//	}
//	return user.Name, user.Avatar, user.BackgroundImage, user.Signature, nil
//}
//
//// 获取用户昵称
//func GetUserNameByID(userID uint) (string, error) {
//	user := model.User{}
//	if err := DB.Select("name").Where("id = ?", userID).First(&user).Error; err != nil {
//		return "", err
//	}
//	return user.Name, nil
//}
//
//// 获取用户点赞的所有视频ID
//func GetLikedVideoIDs(userID uint) ([]uint, error) {
//	var videoIDs []uint
//	err := DB.Model(&model.Like{}).Where("user_id = ? AND liked = 1", userID).Pluck("video_id", &videoIDs).Error
//	if err != nil {
//		return nil, err
//	}
//	return videoIDs, nil
//}
//
//func GetVideoDetailsByID(videoID uint) (uint, string, string, string, error) {
//	video := model.Video{}
//	if err := DB.Select("author_id", "title", "play_url", "cover_url").Where("id = ?", videoID).First(&video).Error; err != nil {
//		return 0, "", "", "", fmt.Errorf("找不到视频ID %d: %v", videoID, err)
//	}
//	return video.AuthorID, video.Title, video.PlayURL, video.CoverURL, nil
//}
//
//// 通过视频ID获取评论总数
//func GetCommentCount(videoID uint) (int64, error) {
//	var count int64
//	if err := DB.Model(&model.Comment{}).Where("video_id = ?", videoID).Where("action_type = ?", 1).Count(&count).Error; err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//
//// 通过视频ID获取点赞总数
//func GetLikeCount(videoID uint) (int64, error) {
//	var count int64
//	if err := DB.Model(&model.Like{}).Where("video_id = ?", videoID).Where("liked = ?", 1).Count(&count).Error; err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//
//// 获取某用户关注总数
//func GetUserFollowCount(userID uint) (int64, error) {
//	var followCount int64
//	if err := DB.Model(&model.Relation{}).Where("user_id = ?", userID).Where("followed = ?", 1).Count(&followCount).Error; err != nil {
//		return 0, err
//	}
//	return followCount, nil
//}
//
//// 获取某用户粉丝总数
//func GetUserFollowerCount(userID uint) (int64, error) {
//	var followerCount int64
//	if err := DB.Model(&model.Relation{}).Where("following_id = ?", userID).Where("followed = ?", 1).Count(&followerCount).Error; err != nil {
//		return 0, err
//	}
//	return followerCount, nil
//}
//
//// 检查一个用户是否关注了另一个用户
//func IsUserFollowingAnotherUser(userID, followingID uint) (bool, error) {
//	var relation model.Relation
//	if err := DB.Where("user_id = ?", userID).Where("following_id = ?", followingID).Where("followed = ?", 1).First(&relation).Error; err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return false, nil // 没有找到关系记录，表示用户没有关注另一个用户
//		}
//		return false, err // 其他错误
//	}
//
//	return true, nil // 找到关系记录，表示用户关注了另一个用户
//}
//
//// 获取用户的获赞总数
//func GetUserTotalReceivedLikes(userID uint) (int64, error) {
//	var totalLikes int64
//	if err := DB.Table("likes").Joins("JOIN videos on videos.id = likes.video_id").Where("videos.author_id = ?", userID).Count(&totalLikes).Error; err != nil {
//		return 0, err
//	}
//	return totalLikes, nil
//}
//
//// 获取用户的点赞总数（用户给出的赞）
//func GetUserTotalGivenLikes(userID uint) (int64, error) {
//	var totalLikes int64
//	if err := DB.Model(&model.Like{}).Where("user_id = ?", userID).Where("liked = ?", 1).Count(&totalLikes).Error; err != nil {
//		return 0, err
//	}
//	return totalLikes, nil
//}
//
//// 获取用户的作品数量
//func GetUserWorkCount(userID uint) (int64, error) {
//	var workCount int64
//	if err := DB.Model(&model.Video{}).Where("author_id = ?", userID).Count(&workCount).Error; err != nil {
//		return 0, err
//	}
//	return workCount, nil
//}
//
//// IsVideoLikedByUser 判断用户是否点赞了视频
//func IsVideoLikedByUser(userID uint, videoID uint) (bool, error) {
//	var count int64
//
//	// 查询 Like 表，看是否有与给定的用户 ID 和视频 ID 匹配的记录
//	if err := DB.Table("likes").Where("user_id = ? AND video_id = ? AND liked = 1", userID, videoID).Count(&count).Error; err != nil {
//		return false, err
//	}
//
//	// 如果 count 大于 0，表示用户点赞了视频
//	return count > 0, nil
//}

//func likeVideo(userID uint, videoID uint) error {
//	// 开始一个新的事务
//	tx := dao.DB.Begin()
//	if tx.Error != nil {
//		return tx.Error
//	}
//
//	// 使用特定的查询构造方式
//	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
//	if err != nil {
//		// 如果记录未找到
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			// 创建一个新的喜欢记录
//			like := model.Like{
//				UserID:  userID,
//				VideoID: videoID,
//				Liked:   1,
//			}
//			// 将新记录保存到数据库
//			if err := tx.Create(&like).Error; err != nil {
//				tx.Rollback() // 回滚事务
//				return err
//			}
//		} else {
//			// 如果发生其他错误，则回滚事务并返回该错误
//			tx.Rollback()
//			return err
//		}
//	}
//
//	// 假设 first 是一个 *model.Like 类型
//	if first.Liked == 1 {
//		tx.Rollback() // 回滚事务
//		return fmt.Errorf("user has already liked this video")
//	}
//
//	// 将喜欢的状态设置为1
//	first.Liked = 1
//	// 保存记录
//	if err := tx.Save(&first).Error; err != nil {
//		tx.Rollback() // 回滚事务
//		return err
//	}
//
//	// 提交事务
//	return tx.Commit().Error
//}
//func unlike(userID uint, videoID uint) error {
//	// 开始一个新的事务
//	tx := dao.DB.Begin()
//	if tx.Error != nil {
//		return tx.Error
//	}
//
//	// 使用特定的查询构造方式
//	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
//	if err != nil {
//		// 如果记录未找到
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			tx.Rollback() // 回滚事务
//			return fmt.Errorf("No like found for this user and video")
//		}
//		// 如果发生其他错误，则回滚事务并返回该错误
//		tx.Rollback()
//		return err
//	}
//
//	// 假设 first 是一个 *model.Like 类型
//	if first.Liked == 0 {
//		tx.Rollback() // 回滚事务
//		return fmt.Errorf("User has already unliked this video")
//	}
//
//	// 将喜欢的状态设置为0
//	first.Liked = 0
//	// 保存记录
//	if err := tx.Save(&first).Error; err != nil {
//		tx.Rollback() // 回滚事务
//		return err
//	}
//
//	// 提交事务
//	return tx.Commit().Error
//}
