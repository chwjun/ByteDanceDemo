package service

import (
	"github.com/RaymondCode/simple-demo/repository"
)

// CommentService 接口定义
type CommentService interface {

	// 给视频模块提供服务
	GetCommentCnt(videoId int64) (int64, error)

	CommentAction(comment repository.Comment) (Comment, error)
	DeleteCommentAction(commentId int64) error
	GetCommentList(videoId int64, userId int64) ([]Comment, error)
}

type Comment struct {
	Id         int64  `json:"id"`
	User       *User  `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
	LikeCount  int64  `json:"like_count"`
	TeaseCount int64  `json:"tease_count"`
}
