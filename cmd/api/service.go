// Package api @Author: youngalone [2023/8/7]
package api

import (
	config2 "bytedancedemo/config"
	"bytedancedemo/dao"
	"bytedancedemo/database/mysql"
	"bytedancedemo/database/redis"
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/router"

	"bytedancedemo/oss"
	"bytedancedemo/utils/log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	config   string
	mode     string
	StartCmd = &cobra.Command{
		Use:   "server",
		Short: "服务入口",
		Long:  "抖音极简版APP服务器",
		PreRun: func(cmd *cobra.Command, args []string) {
			config2.Init(config)
			log.InitLogger(mode) //日志重复
			mysql.Init()
			redis.Init()
			rabbitmq.InitRabbitMQ()
			rabbitmq.InitCommentRabbitMQ()
			rabbitmq.InitFollowRabbitMQ()
			rabbitmq.InitVideoRabbitMQ()
			log.InitLogger(mode) //日志重复
			dao.SetDefault(mysql.DB)
			oss.Init()
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "配置文件路径")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "debug", "运行模式")
}

func run() {
	go router.Setup()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	zap.L().Info("监听中断中...")
	<-quit
	zap.L().Sync()
	zap.L().Info("关闭服务器...")
}
