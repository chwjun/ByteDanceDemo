package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/database"
	"bytedancedemo/model"
	"fmt"
	"sync"
	"time"
)

type VideoServiceImp struct {
	UserService
	//CommentService
	//FavoriteService
}

var (
	size             = 8
	videoServiceImp  *VideoServiceImp // 给controller用的
	videoServiceOnce sync.Once
)

func NewVSIInstance() *VideoServiceImp {
	videoServiceOnce.Do(
		func() {
			videoServiceImp = &VideoServiceImp{
				UserService: &UserServiceImpl{},
				//CommentService:  &CommentServiceImpl{},
				//FavoriteService: &FavoriteServiceImpl{},
			}
		})
	return videoServiceImp
}

func (videoService *VideoServiceImp) Test() {
	fmt.Println("获取接口成功")
}

func (videoService *VideoServiceImp) Feed(latest_time time.Time, user_id int64) ([]ResponseVideo, time.Time, error) {

	// fmt.Println("进入了feed")
	// 根据最新时间查找数据库获取视频的信息
	dao_video_list, err := GetVideosByLatestTime(latest_time)
	fmt.Println(len(dao_video_list))
	if err != nil || len(dao_video_list) == 0 || dao_video_list == nil {
		fmt.Println("Feed")
		return nil, time.Time{}, err
	}
	// fmt.Println("here is Feed")
	response_video_list, err := makeResponseVideo(dao_video_list, videoService, int64(user_id))
	if err != nil {
		return nil, dao_video_list[len(dao_video_list)-1].CreatedAt, err
	}
	return response_video_list, dao_video_list[len(dao_video_list)-1].CreatedAt, nil
}

// 构造返回的视频流
func makeResponseVideo(dao_video_list []*model.Video, videoService *VideoServiceImp, user_id int64) ([]ResponseVideo, error) {
	response_video_list := make([]ResponseVideo, 0, size)
	for _, video := range dao_video_list {
		temp_response_video := ResponseVideo{}
		var wait_group sync.WaitGroup
		wait_group.Add(5)
		//根据作者id查作者信息
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			// author_id := video.AuthorID
			// author, err := videoService.GetUserDetailsById(author_id, &user_id)
			// if err == nil {
			// 	temp_response_video.Author = *author
			// } else {
			// 	return
			// }
			author := User{
				Id: 1,
			}
			temp_response_video.Author = author
			wait_group.Done()
		}(video, &temp_response_video)
		// 根据视频id找评论总数
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			// comment_count, err := videoService.GetCommentCnt(video.ID)
			// if err == nil {
			// 	temp_response_video.Comment_count = int(comment_count)
			// } else {
			// 	return
			// }
			comment_count := 10
			temp_response_video.Comment_count = int64(comment_count)
			wait_group.Done()
		}(video, &temp_response_video)
		// 根据视频id找点赞总数
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			// like_count, err := GetLikeCount(uint(video.ID))
			// if err == nil {
			// 	temp_response_video.Favorite_count = int(like_count)
			// } else {
			// 	return
			// }
			like_count := 100
			temp_response_video.Favorite_count = int64(like_count)
			wait_group.Done()
		}(video, &temp_response_video)
		// 根据当前用户id和视频id判断是否点赞了
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			// is_like, err := IsVideoLikedByUser(uint(user_id), uint(video.ID))s
			// if err == nil {
			// 	temp_response_video.Is_favorite = is_like
			// } else {
			// 	return
			// }
			is_like := true
			temp_response_video.Is_favorite = is_like
			wait_group.Done()
		}(video, &temp_response_video)
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			temp_response_video.Id = int64(video.ID)
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

// 提供给外部的接口
func (videoService *VideoServiceImp) GetVideoListByAuthorID(authorId int64) ([]*model.Video, error) {
	dao_video_list, err := DAOGetVideoListByAuthorID(authorId)
	if err != nil {
		return nil, err
	} else {
		return dao_video_list, nil
	}
}

func (videoService *VideoServiceImp) GetVideoCountByAuthorID(authorId int64) (int, error) {
	dao_video_list, err := DAOGetVideoListByAuthorID(authorId)
	if err != nil {
		return 0, err
	} else {
		return len(dao_video_list), nil
	}
}

func (videoService *VideoServiceImp) PublishList(user_id int64) ([]ResponseVideo, error) {

	dao_video_list, err := videoService.GetVideoListByAuthorID(user_id)
	if err != nil {
		return nil, err
	}
	response_video_list, err := makeResponseVideo(dao_video_list, videoService, user_id)
	if err != nil {
		return nil, err
	}
	return response_video_list, err

}

const Video_list_size = 10

func GetVideosByLatestTime(latest_time time.Time) ([]*model.Video, error) {
	dao.SetDefault(database.DB)
	// 在这里查询
	V := dao.Video
	fmt.Println(V)
	result, err := V.Where(V.CreatedAt.Lt(latest_time)).Order(V.CreatedAt.Desc()).Limit(Video_list_size).Find()
	// fmt.Println(latest_time)
	// fmt.Println(len(result))
	if err != nil {
		fmt.Println("查询最新时间的videos出错了")
		result = nil
		return nil, err
	}
	return result, err
}

func DAOGetVideoListByAuthorID(authorId int64) ([]*model.Video, error) {
	V := dao.Video
	// fmt.Println(V)
	result, err := V.Where(V.AuthorID.Eq(authorId)).Order(V.CreatedAt.Desc()).Find()
	if err != nil || result == nil || len(result) == 0 {
		return nil, err
	}
	return result, err
}
