package service

import (
	"context"
	"fmt"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/proto"
)

type FavoriteServiceImpl struct {
}

const (
	SuccessCode    int32  = 0
	ErrorCode      int32  = 1
	SuccessMessage string = "Success"
)

func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *proto.FavoriteActionRequest) (*proto.FavoriteActionResponse, error) {
	// 使用常量初始化默认值
	statusCode := SuccessCode
	statusMsg := SuccessMessage
	//userID, exists := c.Get("userID")

	userID := uint(1)

	switch *req.ActionType {
	case 1:
		err := dao.Like(userID, uint(*req.VideoId))
		if err != nil {
			statusCode = ErrorCode
			statusMsg = fmt.Sprintf("Failed to like video: %v", err)
		}
	case 2:
		err := dao.Unlike(userID, uint(*req.VideoId))
		if err != nil {
			statusCode = ErrorCode
			statusMsg = fmt.Sprintf("Failed to unlike video: %v", err)
		}
	default:
		err := fmt.Errorf("invalid action_type: %v", *req.ActionType)
		statusCode = ErrorCode
		statusMsg = err.Error()
		return &proto.FavoriteActionResponse{
			StatusCode: &statusCode,
			StatusMsg:  &statusMsg,
		}, err
	}

	return &proto.FavoriteActionResponse{
		StatusCode: &statusCode,
		StatusMsg:  &statusMsg,
	}, nil
}
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *proto.FavoriteListRequest) (*proto.FavoriteListResponse, error) {
	// 从上下文获得userId（假设中间件已将userId放入上下文）
	userID := req.GetUserId()

	// 通过userId获取用户点赞的视频列表
	videoList, err := GetFavoriteVideoInfoByUserID(userID)
	if err != nil {
		errorCode := ErrorCode
		errorMessage := "获取视频失败: " + err.Error()
		return &proto.FavoriteListResponse{
			StatusCode: &errorCode,
			StatusMsg:  &errorMessage,
		}, nil
	}

	successCode := SuccessCode
	successMessage := SuccessMessage
	response := &proto.FavoriteListResponse{
		StatusCode: &successCode,
		StatusMsg:  &successMessage,
		VideoList:  videoList,
	}

	return response, nil
}

func GetFavoriteVideoInfoByUserID(userID int64) ([]*proto.Video, error) {
	videoIDs, err := dao.GetLikedVideoIDs(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("获取点赞视频ID失败: %v", err)
	}

	var protoVideos []*proto.Video
	for _, videoID := range videoIDs {
		videoIDInt64 := int64(videoID)
		videoID := videoID
		//打印videoID
		fmt.Println(videoID)
		authorID, title, playURL, coverURL, err := dao.GetVideoDetailsByID(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取视频详情失败: %v", err)
		}

		commentCount, err := dao.GetCommentCount(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取评论总数失败: %v", err)
		}

		likeCount, err := dao.GetLikeCount(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取点赞总数失败: %v", err)
		}

		isFavorite, err := dao.IsVideoLikedByUser(uint(userID), videoID)
		if err != nil {
			return nil, fmt.Errorf("判断用户是否点赞了视频失败: %v", err)
		}

		requestingUserID := int64(userID)
		author, err := GetUserInfoByID(&requestingUserID, int64(authorID))
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %v", err)
		}

		video := &proto.Video{
			Id:            &videoIDInt64,
			Author:        author,
			PlayUrl:       &playURL,
			CoverUrl:      &coverURL,
			FavoriteCount: &likeCount,
			CommentCount:  &commentCount,
			IsFavorite:    &isFavorite,
			Title:         &title,
		}
		protoVideos = append(protoVideos, video)
	}

	return protoVideos, nil
}

func GetUserInfoByID(requestingUserID *int64, userID int64) (*proto.User, error) {
	name, avatar, backgroundImage, signature, err := dao.GetUserDetailsByID(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取关注总数
	followCount, err := dao.GetUserFollowCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取粉丝总数
	followerCount, err := dao.GetUserFollowerCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 检查是否已关注
	isFollow, err := dao.IsUserFollowingAnotherUser(uint(*requestingUserID), uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取获赞总数
	totalFavorited, err := dao.GetUserTotalReceivedLikes(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取作品数量
	workCount, err := dao.GetUserWorkCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取点赞数量
	favoriteCount, err := dao.GetUserTotalGivenLikes(uint(userID))
	if err != nil {
		return nil, err
	}

	user := &proto.User{
		Id:              &userID,
		Name:            &name,
		FollowCount:     &followCount,
		FollowerCount:   &followerCount,
		IsFollow:        &isFollow,
		Avatar:          &avatar,
		BackgroundImage: &backgroundImage,
		Signature:       &signature,
		TotalFavorited:  &totalFavorited,
		WorkCount:       &workCount,
		FavoriteCount:   &favoriteCount,
	}

	return user, nil
}
