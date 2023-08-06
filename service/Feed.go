package service

import "time"

type ResponseVideo struct {
	id int
	// 作者信息
	play_url       string
	cover_url      string
	favorite_count int
	comment_count  int
	is_favorite    bool
	title          string
}

// 数据库中视频表的结构
type Video_Dao struct {
	// BaseModel
	ID        uint
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt time.Time
	// 视频表
	AuthorID uint
	Title    string
	PlayURL  string
	CoverURL string
}

const size = 10

func Feed(latest_time string) ([]ResponseVideo, string, error) {
	video_dao_list := make([]Video_Dao, 0, size)
	// 根据最新时间查找数据库获取视频的信息
	// 遍历video_list，根据视频id查作者信息
	// 根据视频id找评论总数
	// 根据视频id找点赞总数
	// 根据当前用户id和视频id判断是否点赞了
	// 将上述信息组装成响应的列表
}
