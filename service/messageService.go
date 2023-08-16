package service

import (
	"github.com/RaymondCode/simple-demo/controller"
	"time"
)

type MessageService interface {
	//Send message function
	SendMessage(userId int64, toUserId int64, content string, actionType int64) error

	//Get chat history function
	GetChatHistory(userId int64, toUserId int64, searchTime time.Time) ([]controller.Message, error)

	//user part
	GetLatestMessage(userId int64, selectedUserId int64) (controller.Message, error)
}
