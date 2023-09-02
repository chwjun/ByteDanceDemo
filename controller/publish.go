package controller

import (
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/service"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VideoListResponse struct {
	Response
	VideoList []service.ResponseVideo `json:"video_list"`
}

const temp_pre = "./temp/"

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	data, err := c.FormFile("data")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
	}

	// 以一个不同的字符串作为视频名字，保存到本地tmp文件夹中
	videoName := uuid.New().String()
	if _, err := os.Stat(temp_pre); os.IsNotExist(err) {
		err := os.Mkdir(temp_pre, os.ModePerm)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("[%s : %d] %s \n", file, line, err.Error())
		}
	}
	tempPath := temp_pre + videoName + ".mp4"
	err = c.SaveUploadedFile(data, tempPath)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
	}
	defer os.Remove(tempPath)

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
	err = videoservice.Action(title, user_id, videoName)
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
	log.Println("Publish list !!!!!!!!!!")
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
	err := publishlistMQ.PublishRequest("publishlist")
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
