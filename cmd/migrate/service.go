// Package migrate @Author: youngalone [2023/8/7]
package migrate

import (
	config2 "github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/utils/gen"
	"github.com/gookit/slog"
	"github.com/spf13/cobra"
)

var (
	config   string
	StartCmd = &cobra.Command{
		Use:     "init",
		Short:   "初始化数据库",
		Long:    "init 通过gen生成安全的dao层",
		Example: "go run main.go init -c=\"config/settings.yml\"",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "配置文件路径")
}

func run() {
	slog.Info("开始环境初始化...")
	config2.Init(config)
	database.Init()
	gen.Setup()
	slog.Info("环境初始化成功！")
}
