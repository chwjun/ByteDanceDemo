package controller

import (
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []service.ResponseVideo `json:"video_list,omitempty"`
	NextTime  int64                   `json:"next_time,omitempty"`
}

// 参数latest_time 和 token
func Feed(c *gin.Context) {
	// 传时间戳
	default_time := time.Now().UnixMilli()
	log.Println("Feed!!!!!!!!!!!!!")
	var temp int64
	var latest_time time.Time
	var latest_time_str = c.Query("latest_time")
	if latest_time_str == "" {
		temp = default_time
		latest_time = time.UnixMilli(temp)
	} else {
		temp, err := strconv.ParseInt(latest_time_str, 10, 64)
		if err != nil {
			fmt.Println("%s cannot change to int64", latest_time_str)
			panic(1)
		}
		// 时间只能是比现在小的
		if temp > default_time {
			temp = default_time
		}
		if temp <= 0 {
			temp = time.Now().UnixMilli()
		}
		latest_time = time.UnixMilli(temp)
	}

	// fmt.Println(latest_time)
	videoservice := service.NewVSIInstance()
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

	videoservice.Test()
	// 使用消息队列

	feedMQ := rabbitmq.SimpleVideoFeedMq
	err := feedMQ.PublishRequest("feed")
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "消息队列已满或消息队列出错"},
		})
		return
	}
	// log.Println("运行时间:", time.Now().Sub(t1))
	t2 := time.Now()
	video_list, last_time, err := videoservice.Feed(latest_time, user_id)
	last_time1 := last_time.UnixMilli()
	if err != nil {
		fmt.Println("error happend")
	}
	log.Println("运行时间:", time.Now().Sub(t2))
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: ""},
		VideoList: video_list,
		NextTime:  last_time1,
	})
}
