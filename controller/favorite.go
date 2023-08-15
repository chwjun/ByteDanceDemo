package controller

import (
	"net/http"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type FavoriteActionRequest struct {
	VideoID    int64 `json:"video_id"`
	ActionType int32 `json:"action_type"`
}

func FavoriteAction(c *gin.Context) {
	var req FavoriteActionRequest

	// 尝试将请求数据绑定到我们的请求结构中
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建服务
	s := service.FavoriteServiceImpl{}

	// 调用服务层的FavoriteAction方法
	resp, err := s.FavoriteAction(req.VideoID, req.ActionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有错误，则返回正常响应
	c.JSON(http.StatusOK, resp)
}

func FavoriteList(c *gin.Context) {
	// 从上下文中获取userID
	userID := int64(1)
	//userIDValue, exists := c.Get("userID")

	//if !exists {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
	//	return
	//}

	// 断言userID为期望的类型，例如int64
	//userID, ok := userIDValue.(int64)
	//if !ok {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has an incorrect type"})
	//	return
	//}

	// 创建服务
	s := service.FavoriteServiceImpl{}
	resp, err := s.FavoriteList(userID) // 传递userID给服务层方法
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}

	// 检查服务层返回的状态码，如果不为0，则返回错误
	if resp.StatusCode != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.StatusMsg})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, resp)
}
