package main

import (
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/router"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func main() {
	//go service.RunMessageServer()

	//initialize.ReadConfig()
	//initialize.ConectDB()
	favoriteService := service.NewFavoriteService("43.140.203.85:6379", "", 0)
	go service.StartSyncTask(favoriteService.RedisClient, dao.DB, time.Minute*5) // 每5分钟同步一次 每5分钟同步一次
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("favoriteService", favoriteService)
		c.Next()
	})
	router.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
