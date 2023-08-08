package controller

import (
	"net/http"

	"github.com/RaymondCode/simple-demo/proto"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	var req proto.FavoriteActionRequest

	// 尝试将请求数据绑定到我们的请求结构中
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建服务
	s := service.FavoriteServiceImpl{}

	// 调用服务层的FavoriteAction方法
	resp, err := s.FavoriteAction(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果没有错误，则返回正常响应
	c.JSON(http.StatusOK, resp)
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
