// Package config @Author: youngalone [2023/8/4]
package config

import (
	"github.com/gookit/slog"
	"github.com/spf13/viper"
)

const OSS_USER_AVATAR_DIR = "your_avatar_directory_here"

// 视频模块相关配置
const (
	VIDEO_NUM_PER_REFRESH     = 6
	VIDEO_INIT_NUM_PER_AUTHOR = 10
	// 阿里 OSS 相关配置
	OSS_ACCESS_KEY_ID     = "OSS_ACCESS_KEY_ID"
	OSS_ACCESS_KEY_SECRET = "OSS_ACCESS_KEY_SECRET"
	OSS_BUCKET_NAME       = "OSS_BUCKET_NAME"
	OSS_ENDPOINT          = "OSS_ENDPOINT"
	CUSTOM_DOMAIN         = "CUSTOM_DOMAIN"
	OSS_VIDEO_DIR         = "OSS_VIDEO_DIR"
	PLAY_URL_PREFIX       = CUSTOM_DOMAIN + OSS_VIDEO_DIR
	COVER_URL_SUFFIX      = "?x-oss-process=video/snapshot,t_2000,m_fast"
)

func Init(configFilePath string) {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		slog.Errorf("加载配置信息失败 %v", err)
	} else {
		slog.Debug("加载配置信息成功")
	}
}
