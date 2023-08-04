package main

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()
	config.Init()
	database.Init()
	r := gin.New()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
