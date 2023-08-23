package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"github.com/gookit/slog"
	"time"
)

type MessageServiceImpl struct {
}

var (
	messageServiceImpl *MessageServiceImpl
)

func GetMessageServiceInstance() *MessageServiceImpl {
	once.Do(func() {
		messageServiceImpl = &MessageServiceImpl{}
	})
	return messageServiceImpl
}

func (messageService *MessageServiceImpl) SendMessage(userId int64, toUserId int64, actionType int32, content string) error {
	if actionType != 1 {
		slog.Fatalf("Undefined actionType: %d", actionType)
		return nil
	}
	var message model.Message
	message.SenderID = userId
	message.ReceiverID = toUserId
	message.ActionType = actionType
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	message.Content = content
	m := dao.Message
	result := m.Create(&message)
	if result != nil {
		slog.Fatalf("Send message failed %v")
		return result
	}
	return result
}

func GetChatHistory(userId int64, toUserId int64, lastTime time.Time) ([]*model.Message, error) {
	m := dao.Message
	msg, err := m.Where(m.CreatedAt.Gt(lastTime), m.CreatedAt.Lt(time.Now())).
		Where(m.SenderID.Eq(userId), m.ReceiverID.Eq(toUserId)).Or(m.SenderID.Eq(toUserId), m.ReceiverID.Eq(userId)).
		Order(m.CreatedAt).
		Find()
	if err != nil {
		slog.Fatalf("GetChatHistory failed! %v", err)
	}
	return msg, nil
}

func GetLatestMessage(userId int64, selectedUserId int64) (*model.Message, error) {
	m := dao.Message
	msg, err := m.Where(m.SenderID.Eq(userId), m.ReceiverID.Eq(selectedUserId)).
		Or(m.SenderID.Eq(selectedUserId), m.ReceiverID.Eq(userId)).
		Order(m.CreatedAt.Desc()).
		First()
	if err != nil {
		slog.Fatalf("Fetch latest message failed! %v", err)
	}
	return msg, nil
}
