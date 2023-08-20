package service

import (
	"bytedancedemo/model"
	"time"
)

type ResponseVideo struct {
	Id int64 `json:"id"`
	// 作者信息
	Author         User   `json:"author"`
	Play_url       string `json:"play_url"`
	Cover_url      string `json:"cover_url"`
	Favorite_count int64  `json:"favourite_count"`
	Comment_count  int64  `json:"comment_count"`
	Is_favorite    bool   `json:"is_favourate"`
	Title          string `json:"title"`
}

type VideoService interface {
	// 这里的Video是Feed接口返回的视频列表，不是数据库中的视频列表
	Feed(latest_time time.Time, user_id int64) ([]ResponseVideo, time.Time, error)
	// 这里data的数据类型不太懂
	//Action(data file, title string) (int64, string, error)
	// 这里的Video是list接口返回的视频列表，不是数据库中的视频列表
	PublishList(user_id string) ([]ResponseVideo, error)
	GetVideoListByAuthorID(authorId int64) ([]*model.Video, error)
	GetVideoCountByAuthorID(authorId int64) (int, error)
	Test()
}
