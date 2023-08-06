package database

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;comment:主键"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:记录创建时间"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:记录更新时间"`
	DeletedAt time.Time `gorm:"index;comment:软删除时间"`
}
