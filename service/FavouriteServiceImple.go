package service

import (
	"errors"
	"fmt"

	"bytedancedemo/model"

	"github.com/gookit/slog"
	"gorm.io/gorm"

	"bytedancedemo/dao"
)

type FavoriteServiceImpl struct {
}

func (s *FavoriteServiceImpl) FavoriteAction(videoID int64, actionType int32) (FavoriteActionResponse, error) {
	// 使用常量初始化默认值
	statusCode := SuccessCode
	statusMsg := SuccessMessage

	userID := uint(1)

	switch actionType {
	case 1:
		err := likeVideo(userID, uint(videoID))
		if err != nil {
			statusCode = ErrorCode
			statusMsg = fmt.Sprintf("Failed to like video: %v", err)
		}
	case 2:
		err := unlike(userID, uint(videoID))
		if err != nil {
			statusCode = ErrorCode
			statusMsg = fmt.Sprintf("Failed to unlike video: %v", err)
		}
	default:
		err := fmt.Errorf("invalid action_type: %v", actionType)
		statusCode = ErrorCode
		statusMsg = err.Error()
		return FavoriteActionResponse{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		}, err
	}

	return FavoriteActionResponse{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	}, nil
}

func (s *FavoriteServiceImpl) FavoriteList(userID int64) (FavoriteListResponse, error) {
	// 通过userId获取用户点赞的视频列表
	videoList, err := s.GetFavoriteVideoInfoByUserID(userID)
	if err != nil {
		errorCode := ErrorCode
		errorMessage := "获取视频失败: " + err.Error()
		return FavoriteListResponse{
			StatusCode: errorCode,
			StatusMsg:  errorMessage,
		}, nil
	}

	successCode := SuccessCode
	successMessage := SuccessMessage
	response := FavoriteListResponse{
		StatusCode: successCode,
		StatusMsg:  successMessage,
		VideoList:  videoList,
	}

	return response, nil
}

func (s *FavoriteServiceImpl) GetFavoriteVideoInfoByUserID(userID int64) ([]*ResponseVideo, error) {
	videoIDs, err := GetLikedVideoIDs(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("获取点赞视频ID失败: %v", err)
	}

	var videos []*ResponseVideo
	for _, videoID := range videoIDs {
		// 使用特定的查询构造方式获取视频详情
		authorID, title, playURL, coverURL, err := GetVideoDetailsByID(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取视频详情失败: %v", err)
		}

		commentCount, err := GetCommentCount(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取评论总数失败: %v", err)
		}

		likeCount, err := GetLikeCount(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取点赞总数失败: %v", err)
		}

		isFavorite, err := IsVideoLikedByUser(uint(userID), videoID)
		if err != nil {
			return nil, fmt.Errorf("判断用户是否点赞了视频失败: %v", err)
		}

		requestingUserID := int64(userID)
		author, err := s.GetUserInfoByID(&requestingUserID, int64(authorID))
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %v", err)
		}

		video := &ResponseVideo{
			Id:             int(videoID),
			Author:         *author,
			Play_url:       playURL,
			Cover_url:      coverURL,
			Favorite_count: int(likeCount),
			Comment_count:  int(commentCount),
			Is_favorite:    isFavorite,
			Title:          title,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (s *FavoriteServiceImpl) GetUserInfoByID(requestingUserID *int64, userID int64) (*User, error) {
	// 使用特定的查询构造方式获取用户详情
	name, avatar, backgroundImage, signature, err := GetUserDetailsByID(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取关注总数
	followCount, err := GetUserFollowCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取粉丝总数
	followerCount, err := GetUserFollowerCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 检查是否已关注
	isFollow, err := IsUserFollowingAnotherUser(uint(*requestingUserID), uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取获赞总数
	totalFavorited, err := GetUserTotalReceivedLikes(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取作品数量
	workCount, err := GetUserWorkCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取点赞数量
	favoriteCount, err := GetUserTotalReceivedLikes(uint(userID))
	if err != nil {
		return nil, err
	}

	user := &User{
		Id:              userID,
		Name:            name,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        isFollow,
		Avatar:          avatar,
		BackgroundImage: backgroundImage,
		Signature:       signature,
		TotalFavorited:  totalFavorited,
		WorkCount:       workCount,
		FavoriteCount:   favoriteCount,
	}

	return user, nil
}

func likeVideo(userID int64, videoID int64) error {
	// 使用特定的查询构造方式
	// SELECT * FROM likes WHERE user_id = userID AND video_id = videoID;
	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
	if err != nil {
		// 如果记录未找到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建一个新的喜欢记录
			like := model.Like{
				UserID:  userID,
				VideoID: videoID,
				Liked:   1,
			}
			// 将新记录保存到数据库
			return dao.DB.Create(&like).Error
		} else {
			// 如果发生其他错误，则返回该错误
			return err
		}
	}

	// 假设 first 是一个 *model.Like 类型
	if first.Liked == 1 {
		slog.Error("user has already liked this video")
		return fmt.Errorf("user has already liked this video")
	}

	// 将喜欢的状态设置为1
	first.Liked = 1
	// 保存记录
	return dao.DB.Save(&first).Error
}
func unlike(userID int64, videoID int64) error {
	// 使用特定的查询构造方式
	// SELECT * FROM likes WHERE user_id = userID AND video_id = videoID;
	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
	if err != nil {
		// 如果记录未找到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("No like found for this user and video")
		}
		// 如果发生其他错误，则返回该错误
		return err
	}

	// 假设 first 是一个 *model.Like 类型
	if first.Liked == 0 {
		slog.Error("User has already unliked this video")
		return fmt.Errorf("User has already unliked this video")
	}

	// 将喜欢的状态设置为0
	first.Liked = 0
	// 保存记录
	return dao.DB.Save(&first).Error
}
func GetUserDetailsByID(userID int64) (string, string, string, string, error) {

	first, err := dao.User.Select(dao.User.Name, dao.User.Avatar, dao.User.BackgroundImage, dao.User.Signature).Where(dao.User.ID.Eq(userID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", "", "", fmt.Errorf("No user found for this user ID")
		}
	}
	return first.Name, first.Avatar, first.BackgroundImage, first.Signature, nil
}

func GetLikedVideoIDs(userID int64) ([]uint, error) {

	likes, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.Liked.Eq(1)).Order(dao.Like.CreatedAt.Abs()).Find()

	if err != nil {
		return nil, err
	}

	var videoIDs []uint
	for _, like := range likes {
		videoIDs = append(videoIDs, like.VideoID) // 假设VideoID是model.Like中的一个字段
	}

	return videoIDs, nil
}
func GetVideoDetailsByID(videoID uint) (uint, string, string, string, error) {

	First, err := dao.Video.Select(dao.Video.AuthorID, dao.Video.Title, dao.Video.PlayURL, dao.Video.CoverURL).Where(dao.Video.ID.Eq(videoID)).First()
	if err != nil {
		return 0, "", "", "", fmt.Errorf("找不到视频ID %d: %v", videoID, err)
	}
	return First.AuthorID, First.Title, First.PlayURL, First.CoverURL, nil

}
func GetCommentCount(videoID uint) (int64, error) {

	var count int64
	count, err := dao.Comment.Where(dao.Comment.VideoID.Eq(videoID), dao.Comment.ActionType.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func GetLikeCount(videoID uint) (int64, error) {

	count, err := dao.Like.Where(dao.Like.VideoID.Eq(videoID), dao.Like.Liked.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func GetUserFollowCount(userID uint) (int64, error) {

	count, err := dao.Relation.Where(dao.Relation.UserID.Eq(userID), dao.Relation.Followed.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func GetUserFollowerCount(userID uint) (int64, error) {

	count, err := dao.Relation.Where(dao.Relation.FollowingID.Eq(userID), dao.Relation.Followed.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func IsUserFollowingAnotherUser(userID, followingID uint) (bool, error) {
	relation, err := dao.Relation.Where(dao.Relation.UserID.Eq(userID), dao.Relation.FollowingID.Eq(followingID), dao.Relation.Followed.Eq(1)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // 没有找到关系记录，表示用户没有关注另一个用户
		}
		return false, err // 其他错误
	}

	return relation.Followed == 1, nil // 根据Followed字段返回是否关注
}

func GetUserTotalReceivedLikes(userID uint) (int64, error) {
	likes := dao.Like
	videos := dao.Video

	count, err := likes.Join(videos, videos.ID.EqCol(likes.VideoID)).Where(videos.AuthorID.Eq(userID), likes.Liked.Eq(1)).Count()

	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserWorkCount(userID uint) (int64, error) {

	workCount, err := dao.Video.Where(dao.Video.AuthorID.Eq(userID)).Count()

	if err != nil {
		return 0, err
	}
	return workCount, nil
}

func IsVideoLikedByUser(userID uint, videoID uint) (bool, error) {
	likes := dao.Like

	count, err := likes.Where(likes.UserID.Eq(userID), likes.VideoID.Eq(videoID), likes.Liked.Eq(1)).Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
