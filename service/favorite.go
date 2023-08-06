package service

import (
	"context"
	"fmt"

	"github.com/RaymondCode/simple-demo/proto" // 替换为你的 proto 包路径
	"github.com/gookit/slog"
)

type FavoriteServiceImpl struct {
	// 在这里添加你需要的字段，例如数据库连接，其他服务的引用等
}

func ValidateToken(token *string) bool {

	return true
}

func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *proto.DouyinFavoriteActionRequest) (*proto.DouyinFavoriteActionResponse, error) {

	// 校验token 是否有效
	if !ValidateToken(req.Token) {
		err := fmt.Errorf("invalid token: %v", req.Token)
		slog.Fatalf("%v", err)
		return nil, err
	}

	// 其次，我们需要根据 action_type 来执行相应的操作
	statusCode := int32(0)
	statusMsg := "Success"
	switch req.ActionType {
	case 1:
		// 执行点赞操作
		// 假设我们有一个名为 LikeVideo 的函数来完成这项工作
		// 这只是一个示例，你可能需要替换为你的实际逻辑
		err := LikeVideo(req.VideoId)
		if err != nil {
			statusCode = 1
			statusMsg = fmt.Sprintf("Failed to like video: %v", err)
		}
	case 2:
		// 执行取消点赞操作
		// 假设我们有一个名为 UnlikeVideo 的函数来完成这项工作
		// 这只是一个示例，你可能需要替换为你的实际逻辑
		err := UnlikeVideo(req.VideoId)
		if err != nil {
			statusCode = 1
			statusMsg = fmt.Sprintf("Failed to unlike video: %v", err)
		}
	default:
		// 如果 action_type 不是 1 或者 2，我们就返回一个错误
		return nil, fmt.Errorf("invalid action_type: %v", req.ActionType)
	}

	// 最后，我们返回一个响应
	return &proto.DouyinFavoriteActionResponse{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	}, nil
}

func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *proto.DouyinFavoriteListRequest) (*proto.DouyinFavoriteListResponse, error) {
	// 在这里实现 FavoriteList 方法的逻辑
	// 这只是一个示例，你可能需要替换为你的实际逻辑

	statusCode2 := int32(0)
	return &proto.DouyinFavoriteListResponse{
		StatusCode: &statusCode2,
		VideoList:  []*proto.Video{},
	}, nil
}
