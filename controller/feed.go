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
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// 参数latest_time 和 token
func Feed(c *gin.Context) {
	var video_test = make([]Video, 10)
	default_time := time.Now().UnixMilli()
	var latest_time_str = c.DefaultQuery("latest_time", strconv.FormatInt(default_time, 10))
	temp, err := strconv.ParseInt(latest_time_str, 10, 64)
	if err != nil {
		fmt.Println("%s cannot change to int64", latest_time_str)
		panic(1)
	}
	latest_time := time.UnixMilli(temp)
	// 调用Service的Feed进行处理
	_, _, err = service.Feed(latest_time)
	if err != nil {

	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: video_test,
		NextTime:  time.Now().Unix(),
	})
}
