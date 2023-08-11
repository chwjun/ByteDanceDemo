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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有错误，则返回正常响应
	c.JSON(http.StatusOK, resp)
}

func FavoriteList(c *gin.Context) {
	// 解析请求参数
	var req proto.FavoriteListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters: " + err.Error()})
		return
	}

	// 创建服务
	s := service.FavoriteServiceImpl{}
	resp, err := s.FavoriteList(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	SuccessCode := int32(0)
	// 检查服务层返回的状态码，如果不为0，则返回错误
	if *resp.StatusCode != SuccessCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": *resp.StatusMsg})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, resp)
}
