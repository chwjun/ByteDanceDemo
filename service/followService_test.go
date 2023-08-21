package service

import (
	"bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database"
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/middleware/redis"
	"fmt"
	"log"
	"testing"
)

func TestFollowServiceImp_FollowAction(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	result, err := followServiceImp.FollowAction(1, 2)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}

func TestFollowServiceImp_CancelFollowAction(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	result, err := followServiceImp.CancelFollowAction(1, 2)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}

func TestFollowServiceImp_GetFollowings(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	followings, err := followServiceImp.GetFollowings(1)

	if err != nil {
		log.Default()
	}
	fmt.Println(followings)
}

func TestFollowServiceImp_GetFollowers(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	followers, err := followServiceImp.GetFollowers(2)
	if err != nil {
		log.Default()
	}
	fmt.Println(followers)
}

func TestFollowServiceImp_GetFriends(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	friends, err := followServiceImp.GetFriends(2)
	if err != nil {
		log.Default()
	}
	fmt.Println(friends)
}

func TestFollowServiceImp_GetFollowingCnt(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	userIdCnt, err := followServiceImp.GetFollowingCnt(2)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_GetFollowerCnt(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	userIdCnt, err := followServiceImp.GetFollowerCnt(2)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_CheckIsFollowing(t *testing.T) {
	config.Init("../config/settings.yml")
	database.Init()
	dao.SetDefault(database.DB)
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	result, err := followServiceImp.CheckIsFollowing(1, 2)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}
