package service

import (
	"bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database/mysql"
	"bytedancedemo/database/redis"
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/model"
	"github.com/gookit/goutil/timex"
	"log"
	"testing"
)

func TestGetMessageServiceInstance(t *testing.T) {
	usi1 := GetMessageServiceInstance()
	usi2 := GetMessageServiceInstance()
	if usi1 != usi2 {
		t.Errorf("单例模式出错")
	}
}

func TestMessageServiceImpl_SendMessage(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	redis.Init()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitCommentRabbitMQ()

	var message = model.Message{
		SenderID:   5,
		ReceiverID: 8,
		ActionType: 1,
		Content:    "Test message function",
	}

	err := messageServiceImpl.SendMessage(message.SenderID, message.ReceiverID, message.ActionType, message.Content)
	if err != nil {
		t.Errorf("Send message failed!")
	}
}

func TestMessageServiceImpl_GetChatHistory(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	redis.Init()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitCommentRabbitMQ()

	chat, err := messageServiceImpl.GetChatHistory(5, 8, timex.TodayStart())
	if err != nil {
		t.Errorf("Get chat history failed!")
	}

	for _, msg := range chat {
		log.Println(msg)
	}
}

func TestMessageServiceImpl_GetLatestMessage(t *testing.T) {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
	redis.Init()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitCommentRabbitMQ()

	msg, err := messageServiceImpl.GetLatestMessage(5, 8)
	if err != nil {
		t.Errorf("Get lastest message failed!")
	}
	log.Println(msg)
}

//
