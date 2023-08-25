// Package service @Author: youngalone [2023/8/8]

package service

import "bytedancedemo/model"

type User struct {
	Id              int64  `json:"id"`               // 主键
	Name            string `json:"name"`             // 用户名 用于登录 不可重复
	FollowCount     int64  `json:"follow_count"`     // 关注数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝数
	IsFollow        bool   `json:"is_follow"`        // 是否关注
	Avatar          string `json:"avatar"`           // 头像
	BackgroundImage string `json:"background_image"` // 背景图
	Signature       string `json:"signature"`        // 个性签名
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数
	WorkCount       int64  `json:"work_count"`       // 作品数
	FavoriteCount   int64  `json:"favorite_count"`   // 点赞数
}

type UserService interface {

	// InsertUser 插入用户基础信息
	InsertUser(userBasic *model.User) (*model.User, bool)

	// GetUserBasicByPassword 根据用户名和密码查询基础信息
	GetUserBasicByPassword(username string, password string) (*model.User, bool)

	// GetUserDetailsById 根据用户ID查询获取详细信息
	GetUserDetailsById(id int64, curID *int64) (*User, error)

	// GetUserName 根据用户ID查询用户名
	GetUserName(userId int64) (string, error)
}
