// Package migrate @Author: youngalone [2023/8/7]
package migrate

import (
	config2 "bytedancedemo/config"
	"bytedancedemo/database/mysql"
	"bytedancedemo/database/redis"
	"bytedancedemo/utils/casbin"
	"bytedancedemo/utils/gen"
	"bytedancedemo/utils/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	log2 "log"
)

var (
	config   string
	mode     string
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
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "debug", "运行模式")
}

func run() {
	log2.Println("开始环境初始化...")
	config2.Init(config)
	log.InitLogger(mode)
	mysql.Init()
	redis.Init()
	gen.Setup()
	casbin.Setup()
	zap.L().Info("环境初始化成功！")
}
