// Package mysql @Author: youngalone [2023/8/4]
package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var DSN string

func Init() {
	username := viper.GetString("settings.mysql.username")
	password := viper.GetString("settings.mysql.password")
	host := viper.GetString("settings.mysql.host")
	port := viper.GetString("settings.mysql.port")
	schema := viper.GetString("settings.mysql.schema")
	logLevel := viper.GetInt("settings.mysql.logLevel")
	DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10000ms", username, password, host, port, schema)
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		zap.L().Error("mysql连接失败", zap.String("DSN", DSN), zap.Error(err))
	} else {
		zap.L().Debug("mysql连接成功", zap.String("DSN", DSN))
	}
	DB = db
}
