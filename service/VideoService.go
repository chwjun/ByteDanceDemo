package service

import (
	"bytedancedemo/model"
	"time"
)

type ResponseVideo struct {
	Id int64 `json:"id,omitempty"`
	// 作者信息
	Author         User   `json:"author"`
	Play_url       string `json:"play_url" json:"play_url,omitempty"`
	Cover_url      string `json:"cover_url,omitempty"`
	Favorite_count int64  `json:"favorite_count,omitempty"`
	Comment_count  int64  `json:"comment_count,omitempty"`
	Is_favorite    bool   `json:"is_favorite,omitempty"`
	Title          string `json:"title,omitempt"`
}

type VideoService interface {
	// 这里的Video是Feed接口返回的视频列表，不是数据库中的视频列表
	Feed(latest_time time.Time, user_id int64) ([]ResponseVideo, time.Time, error)
	// 这里data的数据类型不太懂
	//Action(data file, title string) (int64, string, error)
	// 这里的Video是list接口返回的视频列表，不是数据库中的视频列表
	PublishList(user_id int64) ([]ResponseVideo, error)
	GetVideoListByAuthorID(authorId int64) ([]*model.Video, error)
	GetVideoCountByAuthorID(authorId int64) (int64, error)
	Test()
}
