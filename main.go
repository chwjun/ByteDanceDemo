package main

import "github.com/RaymondCode/simple-demo/cmd"

func main() {
	//go service.RunMessageServer()
	//config.Init()
	//utils.Init()
	//r := gin.New()
	//
	//initRouter(r)
	//
	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	cmd.Execute()
}
