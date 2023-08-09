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
	userID, _ := ctx.Value("user_id").(int64)

	// 通过userId获取用户点赞的视频列表
	videoList, err := GetFavoriteVideoInfoByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := &proto.FavoriteListResponse{
		StatusCode: 0,         // 0表示成功
		VideoList:  videoList, // 直接使用从dao返回的视频列表
	}

	return response, nil
}

func GetFavoriteVideoInfoByUserID(userID int64) ([]*proto.Video, error) {
	videoIDs, err := getLikedVideoIDs(userID)
	if err != nil {
		return nil, err
	}

	var protoVideos []*proto.Video
	for _, videoID := range videoIDs {
		video, err := getVideoDetails(videoID, userID)
		if err != nil {
			return nil, err
		}
		protoVideos = append(protoVideos, video)
	}

	return protoVideos, nil
}
