package main

import (
	"bytedancedemo/router"
	"github.com/gin-gonic/gin"
)

func main() {
	//go service.RunMessageServer()

	//initialize.ReadConfig()
	//initialize.ConectDB()
	//go util.StartSyncTask(dao.DB, time.Minute*5)
	r := gin.New()

	router.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
