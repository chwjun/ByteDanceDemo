package initialize

import (
	"fmt"

	"github.com/gookit/slog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Package database @Author: youngalone [2023/8/4]

var DB *gorm.DB

func ConectDB() {
	username := viper.GetString("settings.mysql.username")
	password := viper.GetString("settings.mysql.password")
	host := viper.GetString("settings.mysql.host")
	port := viper.GetString("settings.mysql.port")
	schema := viper.GetString("settings.mysql.schema")
	logLevel := viper.GetInt("settings.mysql.logLevel")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10000ms", username, password, host, port, schema)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		slog.Fatalf("mysql连接失败 %v", err)
	} else {
		slog.Debug("mysql连接成功")
	}
	DB = db
}
