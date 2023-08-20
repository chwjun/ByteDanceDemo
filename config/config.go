// Package config @Author: youngalone [2023/8/4]
package config

import (
	"fmt"

	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

func Init(configFilePath string) {
	// configFilePath = "./settings.yml"
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("出问题了")
		slog.Errorf("加载配置信息失败 %v", err)
	} else {
		slog.Debug("加载配置信息成功")
	}
}
