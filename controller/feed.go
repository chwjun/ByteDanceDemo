package controller

import (
	"net/http"
	"strconv"
	"time"
	service "bytedancedemo/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// 参数latest_time 和 token
func Feed(c *gin.Context) {
	default_time := time.Now().UnixMilli()
	var latest_time = c.DefaultQuery("latest_time",strconv.FormatInt(default_time,10))
	// 调用Service的Feed进行处理
	video_list, next_time, err := service.Feed(latest_time)
	
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}
