package controller

import (
	"bytedancedemo/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ChatResponse struct {
	Response
	MessageList []service.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	/*
		token := c.Query("token")
		toUserId := c.Query("to_user_id")
		content := c.Query("content")

		if user, exist := usersLoginInfo[token]; exist {
			userIdB, _ := strconv.Atoi(toUserId)
			chatKey := genChatKey(user.Id, int64(userIdB))

			atomic.AddInt64(&messageIdSequence, 1)
			curMessage := Message{
				Id:         messageIdSequence,
				Content:    content,
				CreateTime: time.Now().Format(time.Kitchen),
			}

			if messages, exist := tempChat[chatKey]; exist {
				tempChat[chatKey] = append(messages, curMessage)
			} else {
				tempChat[chatKey] = []Message{curMessage}
			}
			c.JSON(http.StatusOK, Response{StatusCode: 0})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		}
	*/
	toUserId := c.Query("to_user_id")
	actionType1 := c.Query("action_type")
	content := c.Query("content")
	userId := c.GetInt64("user_id")
	tarUserId, err := strconv.ParseInt(toUserId, 10, 64)
	actionType, err1 := strconv.ParseInt(actionType1, 10, 32)
	if err != nil || err1 != nil {
		log.Println("toUserId/actionType 参数错误")
		return
	}
	messageService := service.GetMessageServiceInstance()
	err = messageService.SendMessage(userId, tarUserId, int32(actionType), content)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Send Message 接口错误"})
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "消息发送成功"})
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	/*
		token := c.Query("token")
		toUserId := c.Query("to_user_id")

		if user, exist := usersLoginInfo[token]; exist {
			userIdB, _ := strconv.Atoi(toUserId)
			chatKey := genChatKey(user.Id, int64(userIdB))

			c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		}

	*/

	userId := c.GetInt64("user_id")
	toUserId := c.Query("to_user_id")
	preMsgTime := c.Query("pre_msg_time")
	log.Println("preMsgTime", preMsgTime)
	covPreMsgTime, err := strconv.ParseInt(preMsgTime, 10, 64)
	if err != nil {
		log.Println("preMsgTime 参数错误")
		return
	}
	latestTime := time.Unix(covPreMsgTime, 0)
	targetUserId, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		log.Println("toUserId 参数错误")
		return
	}
	messageService := service.GetMessageServiceInstance()
	messages, err := messageService.GetChatHistory(userId, targetUserId, latestTime)
	//log.Println(messages)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "获取消息成功"}, MessageList: service.TransferMsg(messages)})
	}
}
