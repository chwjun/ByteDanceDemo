// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameVideo = "videos"

// Video mapped from table <videos>
type Video struct {
	ID        int64          `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"` // 主键
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime(3);comment:记录创建时间" json:"created_at"`                   // 记录创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime(3);comment:记录更新时间" json:"updated_at"`                   // 记录更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);comment:软删除时间" json:"deleted_at"`                    // 软删除时间
	AuthorID  int64          `gorm:"column:author_id;type:bigint(20) unsigned;comment:视频作者id" json:"author_id"`             // 视频作者id
	Title     string         `gorm:"column:title;type:varchar(191);not null;comment:视频标题" json:"title"`                     // 视频标题
	PlayURL   string         `gorm:"column:play_url;type:varchar(191);not null;comment:视频播放地址" json:"play_url"`             // 视频播放地址
	CoverURL  string         `gorm:"column:cover_url;type:varchar(191);not null;comment:视频封面地址" json:"cover_url"`           // 视频封面地址
}

// TableName Video's table name
func (*Video) TableName() string {
	return TableNameVideo
}