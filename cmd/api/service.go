// Package api @Author: youngalone [2023/8/7]
package api

import (
	config2 "bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database"
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/middleware/redis"
	"bytedancedemo/router"
	//"github.com/RaymondCode/simple-demo/service"
	"github.com/spf13/cobra"
)

var (
	config   string
	StartCmd = &cobra.Command{
		Use:   "server",
		Short: "服务入口",
		Long:  "抖音极简版APP服务器",
		PreRun: func(cmd *cobra.Command, args []string) {
			config2.Init(config)
			database.Init()
<<<<<<< HEAD
			redis.InitRedis()
			rabbitmq.InitRabbitMQ()
			rabbitmq.InitCommentRabbitMQ()
			rabbitmq.InitFollowRabbitMQ()
			//	go service.RunMessageServer()
			dao.SetDefault(database.DB)
=======
			//go service.RunMessageServer()
>>>>>>> d943dc5466637f743705e1147de51792bb031661
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "配置文件路径")
}

func run() {
	router.Setup()
}
