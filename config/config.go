// Package config @Author: youngalone [2023/8/4]
package config

import (
	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile("./config/settings.yml")
	if err := viper.ReadInConfig(); err != nil {
		slog.Errorf("加载配置信息失败 %v", err)
	} else {
		slog.Debug("加载配置信息成功")
	}

}
