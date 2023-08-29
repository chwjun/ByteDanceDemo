package controller

import (
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []service.ResponseVideo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 获取用户id
	user_id := int64(0)
	user_id_temp, exits := c.Get("user_id")

	if !exits {
		user_id = int64(0)
	}
	switch user_id_temp.(type) {
	case int64:
		user_id = user_id_temp.(int64)
	default:
		user_id = int64(0)
	}
	title := c.PostForm("title")
	videoservice := service.NewVSIInstance()
	err = videoservice.Action(data, title, user_id)
	if err != nil {
		log.Println("Action ERROR : ", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "上传或存储失败" + err.Error(),
		})
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  title + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id := int64(0)
	user_id_temp, exits := c.Get("user_id")
	if !exits {
		user_id = int64(0)
	}
	switch user_id_temp.(type) {
	case int64:
		user_id = user_id_temp.(int64)
	default:
		user_id = int64(0)
	}
	publishlistMQ := rabbitmq.SimpleVideoPublishListMq
	err := publishlistMQ.PublishSimpleVideo("publishlist", c)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "消息队列已满或消息队列出错"},
		})
		return
	}
	videoservice := service.NewVSIInstance()
	video_list, err := videoservice.PublishList(user_id)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "[ERROR]:" + err.Error(),
			},
			VideoList: nil,
		})
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: video_list,
	})
}
