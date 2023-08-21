package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/middleware/rabbitmq"

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
// Package service @Author: youngalone [2023/8/18]

// 测试单例模式
func TestGetUserServiceInstance(t *testing.T) {
	usi1 := GetCommentServiceInstance()
	usi2 := GetCommentServiceInstance()
	if usi1 != usi2 {
		t.Errorf("单例模式出错")
	}
}

func TestCommentServiceImpl_CommentAction(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitCommentRabbitMQ()
	//rabbitmq.InitFollowRabbitMQ()
	var comment model.Comment = model.Comment{
		UserID:  5,
		VideoID: 14,
		Content: "这条评论来自单元测试TestInsertComment",
	}
	commentRes, err := commentServiceImpl.CommentAction(comment)
	if err != nil {
		log.Default()
	}
	//fmt.Println(result)
	//redis.InitRedis()
	//var comment model.Comment = model.Comment{
	//	UserID:  5,
	//	VideoID: 14,
	//	Content: "这条评论来自单元测试TestInsertComment",
	//}
	//commentRes, err := commentServiceImpl.CommentAction(comment)
	//if err != nil {
	//	log.Default()
	//}
	fmt.Println(commentRes)
}

func TestCommentServiceImpl_DeleteCommentAction(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitCommentRabbitMQ()
	var comment model.Comment = model.Comment{
		UserID:  1,
		VideoID: 4,
		Content: "这条评论来自单元测试TestInsertComment",
	}
	err := commentServiceImpl.DeleteCommentAction(comment.VideoID)
	if err != nil {
		log.Default()
	}
	//fmt.Println(result)
}

//redis.InitRedis()
//err := commentServiceImpl.DeleteCommentAction(1)
//if err != nil {
//	log.Default()
//}

func TestCommentServiceImpl_GetCommentList(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitCommentRabbitMQ()

	commentList, err := commentServiceImpl.GetCommentList(24, 1)
	if err != nil {
		log.Default()
	}
	fmt.Println(commentList)
	//redis.InitRedis()

}
