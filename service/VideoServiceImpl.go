package service

import (
	"bytedancedemo/database"
	"bytedancedemo/model"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type VideoServiceImp struct {
	VideoService
	UserService
}

var (
	size             = 10
	videoServiceImp  *VideoServiceImp // 给controller用的
	videoServiceOnce sync.Once
)

func NewVSIInstance() *VideoServiceImp {
	videoServiceOnce.Do(
		func() {
			videoServiceImp = &VideoServiceImp{
				UserService: &UserServiceImpl{},
			}
		})
	return videoServiceImp
}

func (videoService *VideoServiceImp) Feed(latest_time time.Time) ([]ResponseVideo, time.Time, error) {
	response_video_list := make([]ResponseVideo, 0, size)
	// 根据最新时间查找数据库获取视频的信息
	dao_video_list, err := database.GetVideosByLatestTime(latest_time)
	if err != nil {
		fmt.Println(err)
		return nil, time.Time{}, err
	}

	// 遍历video_list
	response_video_list, err = makeResponseVideo(dao_video_list, videoService)
	if err != nil {
		return nil, dao_video_list[len(dao_video_list)-1].CreatedAt, err
	}
	return response_video_list, dao_video_list[len(dao_video_list)-1].CreatedAt, nil
}

func makeResponseVideo(dao_video_list []*model.Video, videoService *VideoServiceImp) ([]ResponseVideo, error) {
	response_video_list := make([]ResponseVideo, 0, size)
	for _, video := range dao_video_list {
		temp_response_video := ResponseVideo{}
		var wait_group sync.WaitGroup
		wait_group.Add(5)
		//根据作者id查作者信息
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			author_id := video.AuthorID
			autor, err := videoService.GetUserDetailsById(author_id, &author_id)
			if err == nil {
				temp_response_video.Author = *autor
			} else {
				return
			}
			wait_group.Done()
		}(video, &temp_response_video)
		// 根据视频id找评论总数
		go func(video *model.Video, temp_response_video *ResponseVideo) {

			wait_group.Done()
		}(video, &temp_response_video)
		// 根据视频id找点赞总数
		go func(video *model.Video, temp_response_video *ResponseVideo) {

			wait_group.Done()
		}(video, &temp_response_video)
		// 根据当前用户id和视频id判断是否点赞了
		go func(video *model.Video, temp_response_video *ResponseVideo) {

			wait_group.Done()
		}(video, &temp_response_video)
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			temp_response_video.Id = int(video.ID)
			temp_response_video.Play_url = video.PlayURL
			temp_response_video.Cover_url = video.CoverURL
			temp_response_video.Title = video.Title
			wait_group.Done()
		}(video, &temp_response_video)
		wait_group.Wait()
		response_video_list = append(response_video_list, temp_response_video)
	}
	return response_video_list, nil
}

func (videoService *VideoServiceImp) GetVideoListByAuthorID(authorId int64) ([]*model.Video, error) {
	dao_video_list, err := database.GetVideoListByAuthorID(authorId)
	if err != nil {
		return nil, err
	} else {
		return dao_video_list, nil
	}
}

func (videoService *VideoServiceImp) GetVideoCountByAuthorID(authorId int64) (int, error) {
	dao_video_list, err := database.GetVideoListByAuthorID(authorId)
	if err != nil {
		return 0, err
	} else {
		return len(dao_video_list), nil
	}
}

func (videoService *VideoServiceImp) PublishList(user_id string) ([]ResponseVideo, error) {
	response_video_list := make([]ResponseVideo, 0, size)
	userid, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		return nil, err
	}
	dao_video_list, err := videoService.GetVideoListByAuthorID(userid)
	if err != nil {
		return nil, err
	}
	response_video_list, err = makeResponseVideo(dao_video_list, videoService)
	if err != nil {
		return nil, err
	}
	return response_video_list, err

}
