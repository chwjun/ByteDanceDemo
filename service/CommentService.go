package service

import "bytedancedemo/dao"

type CommentService interface {

	// 给视频模块提供服务
	GetCommentCnt(videoId int64) (int64, error)
	CommentAction(comment dao.Comment) (Comment, error) //参数要看看gen
	DeleteCommentAction(commentId int64) error
	GetCommentList(videoId int64, userId int64) ([]Comment, error)
}
