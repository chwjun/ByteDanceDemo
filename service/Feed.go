package service

import (
	"bytedancedemo/database"
	"bytedancedemo/model"
	"fmt"
	"sync"
	"time"
)

type ResponseVideo struct {
	id int
	// 作者信息
	// Author	用户的结构体
	play_url       string
	cover_url      string
	favorite_count int
	comment_count  int
	is_favorite    bool
	title          string
}

const size = 10

func Feed(latest_time time.Time) ([]ResponseVideo, time.Time, error) {

	response_video_list := make([]ResponseVideo, 0, size)
	// 根据最新时间查找数据库获取视频的信息
	dao_video_list, err := database.GetVideosByLatestTime(latest_time)
	if err != nil {
		fmt.Println(err)
		return response_video_list, time.Time{}, err
	}

	// 遍历video_list
	for _, video := range dao_video_list {
		var wait_group sync.WaitGroup
		wait_group.Add(4)
		//根据视频id查作者信息
		go func(video *model.Video) {

			wait_group.Done()
		}(&video)
		// 根据视频id找评论总数
		go func(video *model.Video) {

			wait_group.Done()
		}(&video)
		// 根据视频id找点赞总数
		go func(video *model.Video) {

			wait_group.Done()
		}(&video)
		// 根据当前用户id和视频id判断是否点赞了
		go func(video *model.Video) {

			wait_group.Done()
		}(&video)
		wait_group.Wait()

	}

	return response_video_list, dao_video_list[len(dao_video_list)-1].CreatedAt, nil
}
