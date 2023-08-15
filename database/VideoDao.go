package database

import (
	"bytedancedemo/model"
	"fmt"
	"time"

	"gorm.io/gen/examples/dal/query"
)

const Video_list_size = 10

func GetVideosByLatestTime(latest_time time.Time) ([]model.Video, error) {
	videos_list := make([]model.Video, Video_list_size)
	// 在这里查询
	V := query.Video
	result := V.Where("CreatedAt < ?", latest_time).Order("CreatedAt desc").Limit(Video_list_size).Find(&videos_list)
	if result.Error != nil {
		fmt.Println("获取视频失败")
	}
	if result.RowsAffected == 0 {
		fmt.Println("找不到视频了")
	}
	return videos_list, result.Error
}
