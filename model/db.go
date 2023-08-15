package model

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;comment:主键"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:记录创建时间"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:记录更新时间"`
	DeletedAt time.Time `gorm:"index;comment:软删除时间"`
}

// 用户表
type User struct {
	BaseModel
	Name            string `gorm:"type:varchar(191);not null;unique;comment:用户名"`
	Password        string `gorm:"size:191;not null;comment:用户密码"`
	Role            string `gorm:"type:varchar(191);not null;comment:角色"`
	Avatar          string `gorm:"type:varchar(191);default:'http://yourserver.com/default_avatar.jpg';comment:用户头像链接"`
	BackgroundImage string `gorm:"type:varchar(191);default:'http://yourserver.com/default_background.jpg';comment:用户个人页顶部大图链接"`
	Signature       string `gorm:"type:text;comment:个人简介"`
}

// 视频表
type Video struct {
	BaseModel
	AuthorID uint   `gorm:"index;not null;comment:视频作者id"`
	Title    string `gorm:"type:varchar(191);not null;comment:视频标题"`
	PlayURL  string `gorm:"type:varchar(191);not null;comment:视频播放地址"`
	CoverURL string `gorm:"type:varchar(191);not null;comment:视频封面地址"`
}

// 评论表
type Comment struct {
	BaseModel
	UserID     uint   `gorm:"index;not null;comment:发布评论的用户id"`
	VideoID    uint   `gorm:"index:idx_video_action;not null;comment:评论视频的id"`
	Content    string `gorm:"type:varchar(191);not null;comment:评论的内容"`
	ActionType int    `gorm:"index:idx_video_action;type:ENUM('1', '2');not null;comment:评论行为，1表示已发布评论，2表示删除评论"`
}

// 点赞表
type Like struct {
	BaseModel
	UserID  uint `gorm:"index:idx_user_video_liked;not null;comment:点赞用户id"`
	VideoID uint `gorm:"index:idx_user_video_liked;not null;comment:点赞视频id"`
	Liked   int  `gorm:"index:idx_user_video_liked;not null;default:1;comment:默认1表示已点赞，0表示未点赞"`
}

// 消息表
type Message struct {
	BaseModel
	SenderID   uint   `gorm:"index;not null;comment:发送message的user id"`
	ReceiverID uint   `gorm:"index;not null;comment:接收message的user id"`
	Content    string `gorm:"type:varchar(191);not null;comment:消息内容"`
	ActionType int    `gorm:"type:ENUM('1', '2');not null;comment:消息行为，1表示发送/2表示撤回"`
}

// 关系表
type Relation struct {
	BaseModel
	UserID      uint `gorm:"index:idx_user_following_followed;not null;comment:用户id"`
	FollowingID uint `gorm:"index:idx_user_following_followed;not null;comment:user id关注的用户id"`
	Followed    int  `gorm:"index:idx_user_following_followed;not null;default:0;comment:默认0表示未关注，1表示已关注"`
}
