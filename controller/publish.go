package controller

import (
	"bytedancedemo/service"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VideoListResponse struct {
	Response
	VideoList []service.ResponseVideo `json:"video_list"`
}

const temp_pre = "./temp/"

func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "获取视频资源失败",
		})
	}
	file, err := data.Open()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "视频打开失败",
		})
	}
	defer file.Close()
	uniqueName := uuid.New().String()
	// 拓展名
	ext := filepath.Ext(data.Filename)
	// 视频名字
	objectName := fmt.Sprintf("%s%s", uniqueName, ext)

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
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  title + " uploaded successfully",
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		videoservice := service.NewVSIInstance()
		err = videoservice.Action(title, user_id, objectName, file)
		if err != nil {
			log.Println("Action ERROR : ", err)
			// c.JSON(http.StatusOK, Response{
			// 	StatusCode: 1,
			// 	StatusMsg:  "上传或存储失败" + err.Error(),
			// })
		}
	}()
	wg.Wait()
	log.Println("上传结束")
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	log.Println("Publish list !!!!!!!!!!")
	user_id := int64(0)
	user_id_temp := c.PostForm("user_id")
	user_id, err := strconv.ParseInt(user_id_temp, 10, 64)
	// publishlistMQ := rabbitmq.SimpleVideoPublishListMq
	// err := publishlistMQ.PublishRequest("publishlist")
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取用户id失败"},
		})
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
