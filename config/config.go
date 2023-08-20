// Package config @Author: youngalone [2023/8/4]
package config

import (
	"github.com/spf13/viper"
	"log"
)

func Init(configFilePath string) {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("加载配置信息失败 %v", err)
	} else {
		log.Println("加载配置信息成功！")
	}
}
