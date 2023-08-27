package controller

import (
	"net/http"
	"strconv"

	"bytedancedemo/service"
	"github.com/gin-gonic/gin"
)

type FavoriteActionRequest struct {
	VideoID    int64 `json:"video_id"`
	ActionType int32 `json:"action_type"`
}

func FavoriteAction(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数异常", "data": err})
		return
	}
	actionType64, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	actionType := int32(actionType64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数异常", "data": err})
		return
	}
	req := FavoriteActionRequest{
		VideoID:    videoId,
		ActionType: actionType,
	}
	//zap.L().Debug("点赞请求参数", zap.Any("req", req))

	s := service.FavoriteServiceImpl{}
	s.StartFavoriteAction()                                                         // 创建 FavoriteServiceImpl 实例
	resp, err := s.FavoriteAction(userIDValue.(int64), req.VideoID, req.ActionType) // 使用实例调用方法
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	// 如果没有错误，则返回正常响应
	c.JSON(http.StatusOK, resp)
}

func FavoriteList(c *gin.Context) {
	// 从上下文中获取userID
	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has an incorrect type"})
		return
	}

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
