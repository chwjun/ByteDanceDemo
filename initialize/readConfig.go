package initialize

import (
	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

func ReadConfig() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		slog.Errorf("加载配置信息失败 %v", err)
	} else {
		slog.Debug("加载配置信息成功")
	}
}
