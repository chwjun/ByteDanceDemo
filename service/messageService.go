package service

import (
	"bytedancedemo/model"
	"time"
)

type MessageService interface {
	//Send message function
	SendMessage(userId int64, toUserId int64, actionType int32, content string) error

	//Get chat history function
	GetChatHistory(userId int64, toUserId int64, lastTime time.Time) ([]*model.Message, error)

	//user part
	GetLatestMessage(userId int64, selectedUserId int64) (*model.Message, error)
}
