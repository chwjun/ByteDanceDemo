package initialize

//
//import (
//	"fmt"
//
//	"github.com/gookit/slog"
//	"github.com/spf13/viper"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//	"gorm.io/gorm/logger"
//
//	"github.com/RaymondCode/simple-demo/global" // 确保你已经导入了global包
//)
//
//// Package database @Author: youngalone [2023/8/4]
//
//func ConectDB() {
//	username := viper.GetString("settings.mysql.username")
//	slog.Debug(username)
//	password := viper.GetString("settings.mysql.password")
//	host := viper.GetString("settings.mysql.host")
//	port := viper.GetString("settings.mysql.port")
//	schema := viper.GetString("settings.mysql.schema")
//	logLevel := viper.GetInt("settings.mysql.logLevel")
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10000ms", username, password, host, port, schema)
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
//	})
//	if err != nil {
//		slog.Fatalf("mysql连接失败 %v", err)
//	} else {
//		slog.Debug("mysql连接成功")
//	}
//
//	global.DB = db // 将局部变量db赋值给全局变量global.DB
//}
