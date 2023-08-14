// Package api @Author: youngalone [2023/8/7]
package api

import (
	"github.com/RaymondCode/simple-demo/api"
	config2 "github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/database/mysql"
	"github.com/RaymondCode/simple-demo/database/redis"
	"github.com/RaymondCode/simple-demo/router"
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
			mysql.Init()
			redis.Init()
			go api.RunMessageServer()
			dao.SetDefault(mysql.DB)
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
