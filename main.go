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
	go service.StartSyncTask(dao.DB, time.Minute*5)
	r := gin.New()

	router.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
