// Package log @Author: youngalone [2023/8/20]
package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Config struct {
	Level      zapcore.Level `json:"level"`       // Level 最低日志等级，DEBUG<INFO<WARN<ERROR<FATAL 例如：info-->收集info等级以上的日志
	FileName   string        `json:"file_name"`   // FileName 日志文件位置
	MaxSize    int           `json:"max_size"`    // MaxSize 进行切割之前，日志文件的最大大小(MB为单位)，默认为100MB
	MaxAge     int           `json:"max_age"`     // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxBackups int           `json:"max_backups"` // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
	Compress   bool          `json:"compress"`    // Compress 日志文件是否压缩
	Mode       string        `json:"mode"`        // Mode 日志模式
}

func loadLogConfig() Config {
	return Config{
		Level:      zapcore.Level(viper.GetInt("settings.log.level")),
		FileName:   viper.GetString("settings.log.path"),
		MaxSize:    viper.GetInt("settings.log.maxSize"),
		MaxAge:     viper.GetInt("settings.log.maxAge"),
		MaxBackups: viper.GetInt("settings.log.maxBackups"),
		Compress:   viper.GetBool("settings.log.compress"),
		Mode:       viper.GetString("settings.log.mode"),
	}
}

func InitLogger(mode string) {
	lCfg := loadLogConfig()
	if mode != lCfg.Mode {
		lCfg.Mode = mode
	}
	writeSyncer := getLogWriter(lCfg.FileName, lCfg.MaxSize, lCfg.MaxBackups, lCfg.MaxAge, lCfg.Compress, lCfg.Mode)
	encoder := getEncoder(lCfg.Mode)

	core := zapcore.NewCore(encoder, writeSyncer, lCfg.Level)
	logger := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)
	zap.L().Debug("日志模块初始化成功", zap.String("mode", mode))

}

// 负责设置 encoding 的日志格式
func getEncoder(mode string) zapcore.Encoder {
	var encodeConfig zapcore.EncoderConfig
	if mode == "debug" {
		encodeConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encodeConfig = zap.NewProductionEncoderConfig()
	}

	// 序列化时间 eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.TimeKey = "time"

	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	if mode == "debug" {
		return zapcore.NewConsoleEncoder(encodeConfig)
	}
	return zapcore.NewJSONEncoder(encodeConfig)
}

// 负责日志写入的位置
func getLogWriter(filename string, maxsize, maxBackup, maxAge int, compress bool, mode string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   compress,  // 是否压缩/归档旧文件
	}
	fileWriteSyncer := zapcore.AddSync(lumberJackLogger)
	// 生产模式下 日志只输出到日志文件 降低日志模块损耗
	if mode == "release" {
		return fileWriteSyncer
	}
	consoleWriteSyncer := zapcore.Lock(os.Stdout)
	return zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)
}
