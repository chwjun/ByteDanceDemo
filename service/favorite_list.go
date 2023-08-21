package service

import (
	"fmt"
	"log"

	"github.com/RaymondCode/simple-demo/util"

	"github.com/RaymondCode/simple-demo/dao"
)

type FavoriteList struct {
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

func (s *FavoriteServiceImpl) GetFavoriteVideoInfoByUserID(userID int64) ([]*Video, error) {
	videoIDs, err := GetLikedVideoIDs(uint(userID))
	log.Printf("videoIDs: %+v", videoIDs)
	if err != nil {
		return nil, fmt.Errorf("获取点赞视频ID失败: %v", err)
	}

	videoDetails, err := GetVideoDetailsByIDs(videoIDs)
	if err != nil {
		return nil, fmt.Errorf("获取视频详情失败: %v", err)
	}

	commentCounts, err := GetCommentCounts(videoIDs)
	if err != nil {
		return nil, fmt.Errorf("获取评论总数失败: %v", err)
	}
	likeCounts, err := util.GetVideosLikes(videoIDs)
	if err != nil {
		return nil, fmt.Errorf("获取点赞总数失败: %v", err)
	}

	likedVideos, err := AreVideosLikedByUser(uint(userID), videoIDs)
	if err != nil {
		return nil, fmt.Errorf("判断用户是否点赞了视频失败: %v", err)
	}

	// 获取作者ID列表
	var authorIDs []uint
	for _, detail := range videoDetails {
		authorIDs = append(authorIDs, detail.AuthorID)
	}

	// 获取作者的详细信息
	authors, err := s.GetUserInfoByIDs(userID, authorIDs)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 创建作者ID到详细信息的映射
	authorDetails := make(map[uint]User)
	for _, author := range authors {
		authorDetails[uint(author.ID)] = *author
	}

	var videos []*Video
	for i, videoID := range videoIDs {
		videoDetail := videoDetails[i]
		commentCount := commentCounts[videoID]
		likeCount := likeCounts[videoID]
		isFavorite := likedVideos[videoID]
		authorDetail := authorDetails[videoDetail.AuthorID]
		//打印videid
		log.Printf("videoID: %+v", videoID)
		log.Printf("videoDetail: %+v", videoDetail)
		log.Printf("videoDetail.PlayURL: %s", videoDetail.PlayURL)
		log.Printf("videoDetail.CoverURL: %s", videoDetail.CoverURL)
		log.Printf("videoDetail.Title: %s", videoDetail.Title)
		log.Printf("commentCount: %d", commentCount)
		log.Printf("likeCount: %d", likeCount)
		log.Printf("isFavorite: %v", isFavorite)
		log.Printf("authorDetail: %+v", authorDetail)

		video := &Video{
			ID:            int64(videoID),
			Author:        authorDetail,
			PlayURL:       videoDetail.PlayURL,
			CoverURL:      videoDetail.CoverURL,
			FavoriteCount: likeCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         videoDetail.Title,
		}
		log.Printf("video: %+v", video)
		videos = append(videos, video)
		log.Printf("videos: %+v", videos) /**/

	}
	log.Printf("运行结束")
	return videos, nil
}

func (s *FavoriteServiceImpl) GetUserInfoByIDs(requestingUserID int64, userIDs []uint) ([]*User, error) {
	// 获取用户详情
	userDetails, err := GetUserDetailsByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	// 获取关注总数
	followingCounts, err := GetUserCounts(userIDs, "following")
	if err != nil {
		return nil, fmt.Errorf("获取关注数量出错: %v", err)
	}

	// 获取粉丝数量
	followerCounts, err := GetUserCounts(userIDs, "fans")
	if err != nil {
		return nil, fmt.Errorf("获取粉丝数量出错: %v", err)
	}

	// 检查是否已关注
	isFollowing, err := IsUserFollowingOtherUsers(uint(requestingUserID), userIDs)
	if err != nil {
		return nil, err
	}

	// 获取获赞总数
	totalFavorited, err := GetUserTotalReceivedLikes(userIDs)
	if err != nil {
		return nil, err
	}
	FavoriteCount, err := util.GetUserFavorites(userIDs)

	if err != nil {
		return nil, err
	}
	// 获取作品数量
	workCounts, err := GetUserWorkCounts(userIDs)
	if err != nil {
		return nil, err
	}
	users := make([]*User, len(userIDs))
	for i, userID := range userIDs {
		detail := userDetails[userID]
		user := &User{
			ID:              int64(userID),
			Name:            detail.Name,
			FollowCount:     followingCounts[userID].Count,
			FollowerCount:   followerCounts[userID].Count,
			IsFollow:        isFollowing[userID],
			Avatar:          detail.Avatar,
			BackgroundImage: detail.BackgroundImage,
			Signature:       detail.Signature,
			TotalFavorited:  totalFavorited[userID],
			WorkCount:       workCounts[userID],
			FavoriteCount:   FavoriteCount[userID], // 与 TotalFavorited 相同
		}
		users[i] = user
	}
	return users, nil
}

type UserDetails struct {
	Name            string
	Avatar          string
	BackgroundImage string
	Signature       string
}

func GetUserDetailsByIDs(userIDs []uint) (map[uint]*UserDetails, error) {
	var userDetails map[uint]*UserDetails = make(map[uint]*UserDetails)

	results, err := dao.User.Select(dao.User.ID, dao.User.Name, dao.User.Avatar, dao.User.BackgroundImage, dao.User.Signature).Where(dao.User.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, fmt.Errorf("无法获取用户详细信息: %v", err)
	}

	for _, result := range results {
		userDetails[result.ID] = &UserDetails{
			Name:            result.Name,
			Avatar:          result.Avatar,
			BackgroundImage: result.BackgroundImage,
			Signature:       result.Signature,
		}
	}

	return userDetails, nil
}

func GetLikedVideoIDs(userID uint) ([]uint, error) {

	likes, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.Liked.Eq(1), dao.Like.DeletedAt.IsNull()).Order(dao.Like.CreatedAt.Abs()).Find()

	if err != nil {
		return nil, err
	}

	var videoIDs []uint
	for _, like := range likes {
		videoIDs = append(videoIDs, like.VideoID) // 假设VideoID是model.Like中的一个字段
	}

	return videoIDs, nil
}

type VideoDetail struct {
	AuthorID uint
	Title    string
	PlayURL  string
	CoverURL string
}

// 这个函数的作用是根据视频ID获取视频的详细信息
func GetVideoDetailsByIDs(videoIDs []uint) ([]*VideoDetail, error) {
	var details []*VideoDetail

	results, err := dao.Video.Select(dao.Video.AuthorID, dao.Video.Title, dao.Video.PlayURL, dao.Video.CoverURL).Where(dao.Video.ID.In(videoIDs...)).Find()

	if err != nil {
		return nil, fmt.Errorf("无法获取视频的详细信息: %v", err)
	}

	for _, result := range results {
		detail := &VideoDetail{
			AuthorID: result.AuthorID,
			Title:    result.Title,
			PlayURL:  result.PlayURL,
			CoverURL: result.CoverURL,
		}
		details = append(details, detail)
	}

	return details, nil
}

type VideoCommentCount struct {
	VideoID uint
	Count   int64
}

func GetCommentCounts(videoIDs []uint) (map[uint]int64, error) {
	var counts []*VideoCommentCount

	// 使用IN操作符一次性获取所有视频的评论数量
	err := dao.Comment.
		Where(dao.Comment.VideoID.In(videoIDs...)).
		Where(dao.Comment.ActionType.Eq(1)).
		Select(dao.Comment.VideoID, dao.Comment.ID.Count().As("count")).
		Group(dao.Comment.VideoID).
		Scan(&counts) // 注意这里添加了.Error，用于获取错误信息

	if err != nil {
		return nil, fmt.Errorf("无法获取视频的评论数量: %v", err) // 修改的地方
	}

	// 创建一个映射以快速查找结果
	resultsMap := make(map[uint]int64)
	for _, count := range counts {
		resultsMap[count.VideoID] = count.Count
	}

	// 确保所有传入的视频ID都包含在结果中
	for _, videoID := range videoIDs {
		if _, exists := resultsMap[videoID]; !exists {
			resultsMap[videoID] = 0 // 如果没有评论，则计数为0
		}
	}

	return resultsMap, nil
}

type VideoLikeCount struct {
	VideoID uint
	Count   int64
}

//func GetLikeCounts(videoIDs []uint) (map[uint]int64, error) {
//	var counts []*VideoLikeCount
//
//	// 使用IN操作符一次性获取所有视频的喜欢（like）数量
//	err := dao.Like.Select(dao.Like.VideoID, dao.Like.ID.Count().As("count")).
//		Where(dao.Like.VideoID.In(videoIDs...), dao.Like.Liked.Eq(1), dao.Like.DeletedAt.IsNull()).
//		Group(dao.Like.VideoID).Scan(&counts)
//
//	if err != nil {
//		return nil, fmt.Errorf("无法获取视频的喜欢（like）数量: %v", err)
//	}
//
//	// 创建一个映射以快速查找结果
//	resultsMap := make(map[uint]int64)
//	for _, count := range counts {
//		resultsMap[count.VideoID] = count.Count
//	}
//
//	// 确保所有传入的视频ID都包含在结果中
//	for _, videoID := range videoIDs {
//		if _, exists := resultsMap[videoID]; !exists {
//			resultsMap[videoID] = 0 // 如果没有喜欢（like），则计数为0
//		}
//	}
//
//	return resultsMap, nil
//}

type UserCount struct {
	UserID uint `gorm:"column:user_id"`
	Count  int64
}

func GetUserCounts(userIDs []uint, countType string) (map[uint]*UserCount, error) {
	var counts []*UserCount
	var err error

	if countType == "following" {
		err = dao.Relation.Select(dao.Relation.UserID, dao.Relation.ID.Count().As("count")).
			Where(dao.Relation.UserID.In(userIDs...), dao.Relation.Followed.Eq(1), dao.Relation.DeletedAt.IsNull()).
			Group(dao.Relation.UserID).
			Scan(&counts)
	} else if countType == "fans" {
		err = dao.Relation.Select(dao.Relation.FollowingID.As("user_id"), dao.Relation.ID.Count().As("count")).
			Where(dao.Relation.FollowingID.In(userIDs...), dao.Relation.Followed.Eq(1), dao.Relation.DeletedAt.IsNull()).
			Group(dao.Relation.FollowingID).
			Scan(&counts)
	} else {
		return nil, fmt.Errorf("无效的计数类型: %s", countType)
	}

	if err != nil {
		return nil, fmt.Errorf("无法获取用户的%s数量: %v", countType, err)
	}

	// 创建一个映射以快速查找结果
	resultsMap := make(map[uint]*UserCount)
	for _, count := range counts {
		resultsMap[count.UserID] = count
	}

	// 确保所有传入的用户ID都包含在结果中
	for _, userID := range userIDs {
		if _, exists := resultsMap[userID]; !exists {
			resultsMap[userID] = &UserCount{UserID: userID, Count: 0} // 如果没有粉丝，则计数为0
		}
	}
	return resultsMap, nil
}

type FollowingResult struct {
	FollowingID uint
}

func IsUserFollowingOtherUsers(userID uint, followingIDs []uint) (map[uint]bool, error) {
	followingMap := make(map[uint]bool)

	// 获取当前用户关注的所有用户
	var results []FollowingResult
	err := dao.Relation.Select(dao.Relation.FollowingID).Where(dao.Relation.UserID.Eq(userID), dao.Relation.FollowingID.In(followingIDs...), dao.Relation.Followed.Eq(1), dao.Relation.DeletedAt.IsNull()).Scan(&results)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户关注的人: %v", err)
	}

	// 将所有关注的用户ID添加到map中，并设置值为true
	for _, result := range results {
		followingMap[result.FollowingID] = true
	}

	// 对于用户没有关注的用户，将他们添加到map中，并设置值为false
	for _, id := range followingIDs {
		if _, ok := followingMap[id]; !ok {
			followingMap[id] = false
		}
	}

	return followingMap, nil
}

type LikeResult struct {
	AuthorID uint
	Count    int64
}

func GetUserTotalReceivedLikes(userIDs []uint) (map[uint]int64, error) {
	likesCount := make(map[uint]int64)

	var results []LikeResult
	err := dao.Like.
		Join(dao.Video, dao.Video.ID.EqCol(dao.Like.VideoID)).
		Where(dao.Video.AuthorID.In(userIDs...), dao.Like.Liked.Eq(1), dao.Video.DeletedAt.IsNull(), dao.Like.DeletedAt.IsNull()).
		Group(dao.Video.AuthorID).
		Select(dao.Video.AuthorID, dao.Like.ID.Count().As("count")).
		Scan(&results)

	if err != nil {
		return nil, fmt.Errorf("无法获取用户总接收的喜欢数量: %v", err)
	}

	for _, result := range results {
		likesCount[result.AuthorID] = result.Count
	}

	// 对于没有获取到喜欢数量的用户，将他们添加到map中，并设置值为0
	for _, id := range userIDs {
		if _, ok := likesCount[id]; !ok {
			likesCount[id] = 0
		}
	}

	return likesCount, nil
}

//因为。

//func GetUserTotalReceivedLikes(userIDs []uint) (map[uint]int64, error) {
//	// 初始化一个用户ID到其接收的总喜欢数的映射
//	userTotalLikes := make(map[uint]int64)
//
//	// 遍历每个用户ID
//	for _, userID := range userIDs {
//		// 使用 GetLikedVideoIDs 函数获取该用户喜欢的视频ID
//		videoIDs, err := GetLikedVideoIDs(userID)
//		if err != nil {
//			return nil, fmt.Errorf("无法获取用户（ID: %d）喜欢的视频ID: %v", userID, err)
//		}
//
//		// 使用 GetTotalLikeCounts 函数获取这些视频的总喜欢数
//		totalLikeCount, err := util.GlobalRedisClient.GetTotalLikeCounts(videoIDs)
//		if err != nil {
//
//			return nil, fmt.Errorf("无法获取视频的总喜欢数: %v", err)
//		}
//
//		// 将总喜欢数加到映射中
//		userTotalLikes[userID] = totalLikeCount
//	}
//
//	return userTotalLikes, nil
//}

type UserWorkCount struct {
	AuthorID uint
	Count    int64
}

func GetUserWorkCounts(userIDs []uint) (map[uint]int64, error) {
	workCounts := make(map[uint]int64)

	// 查询结果结构体数组
	var results []UserWorkCount
	err := dao.Video.Where(dao.Video.AuthorID.In(userIDs...), dao.Video.DeletedAt.IsNull()).Group(dao.Video.AuthorID).Select(dao.Video.AuthorID, dao.Video.ID.Count().As("count")).Scan(&results)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户的作品数量: %v", err)
	}

	// 将结果添加到map中
	for _, result := range results {
		workCounts[result.AuthorID] = result.Count
	}

	// 对于没有获取到作品数量的用户，将他们添加到map中，并设置值为0
	for _, id := range userIDs {
		if _, ok := workCounts[id]; !ok {
			workCounts[id] = 0
		}
	}

	return workCounts, nil
}

type LikeVideoResult struct {
	VideoID uint
}

func AreVideosLikedByUser(userID uint, videoIDs []uint) (map[uint]bool, error) {
	likedVideos := make(map[uint]bool)

	var results []LikeVideoResult
	err := dao.Like.Select(dao.Like.VideoID).Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.In(videoIDs...), dao.Like.Liked.Eq(1), dao.Like.DeletedAt.IsNull()).Scan(&results)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户喜欢的视频: %v", err)
	}

	for _, result := range results {
		likedVideos[result.VideoID] = true
	}

	// Set false for videos not liked by the user
	for _, id := range videoIDs {
		if _, ok := likedVideos[id]; !ok {
			likedVideos[id] = false
		}
	}

	return likedVideos, nil
}
