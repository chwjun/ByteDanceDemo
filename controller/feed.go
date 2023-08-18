package controller

import (
	"bytedancedemo/service"
	"fmt"
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
	default_time := time.Now().UnixMilli()
	var latest_time_str = c.DefaultQuery("latest_time", strconv.FormatInt(default_time, 10))
	temp, err := strconv.ParseInt(latest_time_str, 10, 64)
	if err != nil {
		fmt.Println("%s cannot change to int64", latest_time_str)
		panic(1)
	}
	latest_time := time.UnixMilli(temp)
	videoservice := service.NewVSIInstance()
	// 调用Service的Feed进行处理
	user_id := c.GetInt("userID")
	fmt.Println(latest_time)
	fmt.Println(user_id)
	fmt.Println(videoservice)
	video_list, last_time, err := videoservice.Feed(latest_time, user_id)
	last_time1 := last_time.UnixMilli()
	if err != nil {
		fmt.Println("error happend")
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: video_list,
		NextTime:  last_time1,
	})
}
