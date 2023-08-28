package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"fmt"
	"log"
	"sync"
	"time"
)

type VideoServiceImp struct {
	UserService
	CommentService
	FavoriteService
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
				UserService:     &UserServiceImpl{},
				CommentService:  &CommentServiceImpl{},
				FavoriteService: &FavoriteServiceImpl{},
			}
		})
	return videoServiceImp
}

func (videoService *VideoServiceImp) Test() {
	fmt.Println("获取接口成功")
}

func (videoService *VideoServiceImp) Feed(latest_time time.Time, user_id int64) ([]ResponseVideo, time.Time, error) {
	// dao.SetDefault(database.DB)
	// fmt.Println("进入了feed")
	// 根据最新时间查找数据库获取视频的信息
	dao_video_list, err := GetVideosByLatestTime(latest_time)
	// fmt.Println(len(dao_video_list))
	if err != nil || len(dao_video_list) == 0 || dao_video_list == nil {
		fmt.Println("Feed")
		return nil, time.Time{}, err
	}
	// log.Println("数据库操作运行时间:", time.Now().Sub(t1))
	t2 := time.Now()
	// 获取剩余信息，构造返回的结构体
	response_video_list, err := makeResponseVideo(dao_video_list, videoService, int64(user_id))
	if err != nil {
		return nil, dao_video_list[len(dao_video_list)-1].CreatedAt, err
	}
	log.Println("构造运行时间:", time.Since(t2))
	return response_video_list, dao_video_list[len(dao_video_list)-1].CreatedAt, nil
}

// 构造返回的视频流
func makeResponseVideo(dao_video_list []*model.Video, videoService *VideoServiceImp, user_id int64) ([]ResponseVideo, error) {
	// 返回的视频流
	response_video_list := make([]ResponseVideo, len(dao_video_list))
	for _, video := range dao_video_list {
		temp_response_video := ResponseVideo{}
		var wait_group sync.WaitGroup
		wait_group.Add(5)
		t_total := time.Now()
		//根据作者id查作者信息
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			author_id := video.AuthorID
			t1 := time.Now()
			author, err := videoService.GetUserDetailsById(author_id, &user_id)
			log.Println("获取作者信息运行时间:", time.Since(t1))
			if err == nil {
				temp_response_video.Author = *author
			} else {
				temp_response_video.Author = User{}
			}
			// author := User{
			// 	Id: 1,
			// }
			// temp_response_video.Author = author
			wait_group.Done()
		}(video, &temp_response_video)

		// 根据视频id找评论总数
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			t1 := time.Now()
			comment_count, err := videoService.GetCommentCnt(video.ID)
			log.Println("获取评论总数运行时间:", time.Since(t1))
			if err == nil {
				temp_response_video.Comment_count = comment_count
			} else {
				temp_response_video.Comment_count = int64(0)
			}
			// comment_count := 10
			// temp_response_video.Comment_count = comment_count
			wait_group.Done()
		}(video, &temp_response_video)

		// 根据视频id找点赞总数
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			t1 := time.Now()
			like_counts, err := videoService.GetVideosLikes([]int64{video.ID})
			log.Println("获取点赞总数运行时间:", time.Since(t1))
			like_count := like_counts[video.ID]
			if err == nil {
				temp_response_video.Favorite_count = like_count
			} else {
				temp_response_video.Favorite_count = int64(0)
			}
			// like_count := 100
			temp_response_video.Favorite_count = like_count
			wait_group.Done()
		}(video, &temp_response_video)

		// 根据当前用户id和视频id判断是否点赞了
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			t1 := time.Now()
			is_likes, err := videoService.AreVideosLikedByUser(user_id, []int64{video.ID})
			log.Println("判断是否点赞运行时间:", time.Since(t1))
			is_like := is_likes[video.ID]
			if err == nil {
				temp_response_video.Is_favorite = is_like
			} else {
				temp_response_video.Is_favorite = false
			}
			// is_like := true
			wait_group.Done()
		}(video, &temp_response_video)
		// 添加剩余信息
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			t1 := time.Now()
			temp_response_video.Id = video.ID
			temp_response_video.Play_url = video.PlayURL
			temp_response_video.Cover_url = video.CoverURL
			temp_response_video.Title = video.Title
			log.Println("其他运行时间:", time.Since(t1))
			wait_group.Done()
		}(video, &temp_response_video)
		wait_group.Wait()
		log.Println("运行时间:", time.Since(t_total))
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

// 提供给外部的接口
func (videoService *VideoServiceImp) GetVideoCountByAuthorID(authorId int64) (int64, error) {
	dao_video_list, err := DAOGetVideoListByAuthorID(authorId)
	if err != nil {
		return 0, err
	} else {
		return int64(len(dao_video_list)), nil
	}
}

func (videoService *VideoServiceImp) PublishList(user_id int64) ([]ResponseVideo, error) {
	// dao.SetDefault(database.DB)
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

func DAOGetVideoListByAuthorID(authorId int64) ([]*model.Video, error) {
	V := dao.Video
	// fmt.Println(V)
	result, err := V.Where(V.AuthorID.Eq(authorId)).Order(V.CreatedAt.Desc()).Find()
	if err != nil || result == nil || len(result) == 0 {
		return nil, err
	}
	return result, err
}

// 这个是video专用的通过时间获取videolist
func GetVideosByLatestTime(latest_time time.Time) ([]*model.Video, error) {
	// dao.SetDefault(mysql.DB)
	// 在这里查询
	V := dao.Video
	// fmt.Println(V)
	result, err := V.Where(V.CreatedAt.Lt(latest_time)).Order(V.CreatedAt.Desc()).Limit(size).Find()
	// fmt.Println(latest_time)
	// fmt.Println(len(result))
	if err != nil {
		fmt.Println("查询最新时间的videos出错了")
		result = nil
		return nil, err
	}
	return result, err
}
