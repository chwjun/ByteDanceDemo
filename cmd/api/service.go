// Package api @Author: youngalone [2023/8/7]
package api

import (
	config2 "bytedancedemo/config"
	"bytedancedemo/database"
	"bytedancedemo/router"
	"bytedancedemo/service"
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
			go service.RunMessageServer()
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
