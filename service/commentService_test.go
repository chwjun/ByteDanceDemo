package service

import (
	"fmt"
	//"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/middleware/redis"
	"github.com/RaymondCode/simple-demo/model"
	"log"
	"testing"
)

//func TestCommentServiceImpl_GetCommentCnt(t *testing.T) {
//	redis.InitRedis()
//	count, err := commentServiceImpl.GetCommentCnt(25)
//	if err != nil {
//		log.Default()
//	}
//	fmt.Println(count)
//}

func TestCommentServiceImpl_CommentAction(t *testing.T) {
	redis.InitRedis()
	var comment model.Comment = model.Comment{
		UserID:  5,
		VideoID: 14,
		Content: "这条评论来自单元测试TestInsertComment",
	}
	commentRes, err := commentServiceImpl.CommentAction(comment)
	if err != nil {
		log.Default()
	}
	fmt.Println(commentRes)
}

func TestCommentServiceImpl_DeleteCommentAction(t *testing.T) {
	redis.InitRedis()
	err := commentServiceImpl.DeleteCommentAction(1)
	if err != nil {
		log.Default()
	}
}

func TestCommentServiceImpl_GetCommentList(t *testing.T) {
	redis.InitRedis()
	commentList, err := commentServiceImpl.GetCommentList(24, 1)
	if err != nil {
		log.Default()
	}
	fmt.Println(commentList)
}
