package dao

import (
	"github.com/gookit/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	// 配置 MySQL
	//username := viper.GetString("settings.mysql.username")
	//slog.Debug(username)
	//password := viper.GetString("settings.mysql.password")
	//host := viper.GetString("settings.mysql.host")
	//port := viper.GetString("settings.mysql.port")
	//schema := viper.GetString("settings.mysql.schema")
	//logLevel := viper.GetInt("settings.mysql.logLevel")
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10000ms", username, password, host, port, schema)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	//})

	dsn := "sample_douyin:sample_douyin@tcp(43.140.203.85:3306)/sample_douyin?charset=utf8&parseTime=True&loc=Local&timeout=10000ms"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(1)),
	})

	if err != nil {
		slog.Fatalf("mysql连接失败 %v", err)
	} else {
		slog.Debug("mysql连接成功")
	}

	DB = db
}
