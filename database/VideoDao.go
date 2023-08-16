package database

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"time"
)

const Video_list_size = 10

func GetVideosByLatestTime(latest_time time.Time) ([]*model.Video, error) {
	// 在这里查询
	V := dao.Video
	result, err := V.Where(V.CreatedAt.Lt(latest_time)).Order(V.CreatedAt.Desc()).Limit(Video_list_size).Find()
	//result := DB.Where("CreatedAt < ?", latest_time).Order("CreatedAt desc").Limit(Video_list_size).Find(&videos_list)
	if err != nil {
		result = nil
		return nil, err
	}
	return result, err
}

func GetVideoListByAuthorID(authorId int64) ([]*model.Video, error) {
	V := dao.Video
	result, err := V.Where(V.AuthorID.Eq(authorId)).Order(V.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return result, err
}
