package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"errors"
	"github.com/gookit/slog"
	"gorm.io/gen/field"
	"strconv"
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

func (messageService *MessageServiceImpl) GetChatHistory(userId int64, toUserId int64, lastTime time.Time) ([]*model.Message, error) {
	m := dao.Message
	msg, err := m.Where(m.Columns(m.SenderID, m.ReceiverID).In(field.Values([][]interface{}{{userId, toUserId}, {toUserId, userId}}))).
		Where(m.CreatedAt.Gt(lastTime)).
		Order(m.CreatedAt).
		Find()
	if err != nil {
		slog.Fatalf("GetChatHistory failed! %v", err)
	}
	return msg, nil
}

func (messageService *MessageServiceImpl) GetLatestMessage(userId int64, selectedUserId int64) (*model.Message, error) {
	m := dao.Message
	msg, err := m.Where(m.SenderID.Eq(userId), m.ReceiverID.Eq(selectedUserId)).
		Or(m.SenderID.Eq(selectedUserId), m.ReceiverID.Eq(userId)).
		Order(m.CreatedAt.Desc()).
		Find()
	if err != nil {
		slog.Fatalf("Fetch latest message failed! %v", err)
	}
	if len(msg) == 0 {
		return nil, errors.New("没有消息")
	}
	return msg[0], nil
}

func TransferMsg(msg []*model.Message) []Message {
	newMsg := make([]Message, len(msg))
	for i, m := range msg {
		newMsg[i].Id = m.ID
		newMsg[i].ToUserId = m.ReceiverID
		newMsg[i].FromUserId = m.SenderID
		newMsg[i].Content = m.Content
		newMsg[i].CreateTime = strconv.FormatInt(m.CreatedAt.UnixMilli(), 10)
	}
	return newMsg
}
