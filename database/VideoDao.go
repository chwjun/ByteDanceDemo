package database

import (
	"fmt"
	"time"
)

type Video struct {
	BaseModel
	AuthorID uint   `gorm:"index:not null;comment:视频作者id"`
	Title    string `gorm:"type:varchar(255);not null;comment:视频标题"`
	PlayURL  string `gorm:"type:varchar(255);not null;comment:视频播放地址"`
	CoverURL string `gorm:"type:varchar(255);not null;comment:视频封面地址"`
}
const Video_list_size = 10

func (Video) TableName() string {
	return "video"
}

func GetVideosByLatestTime(latest_time time.Time) ([]Video, error) {
	videos_list := make([]Video,Video_list_size)
	result := DB.Where("CreatedAt < ?", latest_time).Order("CreatedAt desc").Limit(Video_list_size).Find(&videos_list)
	if result.Error != nil {
		fmt.Println("获取视频失败")
	}
	if result.RowsAffected == 0{
		fmt.Println("找不到视频了")
	}
	return videos_list, result.Error
}
