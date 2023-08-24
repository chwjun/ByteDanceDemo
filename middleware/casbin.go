// Package middleware @Author: youngalone [2023/8/9]
package middleware

import (
	"bytedancedemo/utils/casbin"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func CasbinMiddleware(c *gin.Context) {
	e, err := casbin.GetCasbin()
	if err != nil {
		zap.L().Error("casbin中间件出错", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("casbin中间件出错 %v", err),
		})
		c.Abort()
		return
	}
	role := c.GetString("role")
	allow, err := e.Enforce(role, c.Request.URL.Path, c.Request.Method)
	if err != nil {
		zap.L().Error("casbin中间件出错", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("casbin中间件出错 %v", err),
		})
		c.Abort()
		return
	}
	if allow {
		c.Next()
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": fmt.Sprintf("对不起 %v 没有 %v-%v 访问权限，请联系管理员", role, c.Request.URL.Path, c.Request.Method),
		})
		c.Abort()
		return
	}
}
