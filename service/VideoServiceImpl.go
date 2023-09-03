package service

import (
	"bytedancedemo/dao"
	"bytedancedemo/middleware/rabbitmq"
	"bytedancedemo/model"
	oss1 "bytedancedemo/oss"
	"fmt"
	"log"
	"mime/multipart"
	"runtime"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	// "github.com/rs/xid"
)

type VideoServiceImp struct {
	UserService
	CommentService
	FavoriteService
}

const Video_list_size = 10
const temp_pre = "./temp/"

var (
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
	zap.L().Debug("Feed获取服务层接口成功")
}

func (videoService *VideoServiceImp) Feed(latest_time int64, user_id int64) ([]ResponseVideo, time.Time, error) {
	// 通过rabbitmq 获取数据库中的数据
	feedMQ := rabbitmq.SimpleVideoFeedMq
	dao_video_list, err := feedMQ.PublishRequest("feed", latest_time, user_id)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		return nil, time.Time{}, err
	}
	// 根据最新时间查找数据库获取视频的信息
	// fmt.Println(len(dao_video_list))
	if err != nil || len(dao_video_list) == 0 || dao_video_list == nil {
		//fmt.Println("Feed")
		return nil, time.Time{}, err
	}
	// 获取剩余信息，构造返回的结构体
	response_video_list, err := makeResponseVideo(dao_video_list, videoService, user_id)
	if err != nil {
		return nil, dao_video_list[len(dao_video_list)-1].CreatedAt, err
	}
	// log.Println("构造运行时间:", time.Since(t2))
	return response_video_list, dao_video_list[len(dao_video_list)-1].CreatedAt, nil
}

// 构造返回的视频流
func makeResponseVideo(dao_video_list []*model.Video, videoService *VideoServiceImp, user_id int64) ([]ResponseVideo, error) {
	// 返回的视频流
	response_video_list := make([]ResponseVideo, len(dao_video_list))
	for index, video := range dao_video_list {
		temp_response_video := ResponseVideo{}
		var wait_group sync.WaitGroup
		wait_group.Add(5)
		// t_total := time.Now()
		//根据作者id查作者信息
		go func(video *model.Video, temp_response_video *ResponseVideo) {
			author_id := video.AuthorID
			// t1 := time.Now()
			author, err := videoService.GetUserDetailsById(author_id, &user_id)
			// log.Println("获取作者信息运行时间:", time.Since(t1))
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
			// t1 := time.Now()
			comment_count, err := videoService.GetCommentCnt(video.ID)
			// log.Println("获取评论总数运行时间:", time.Since(t1))
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
			// t1 := time.Now()
			like_counts, err := videoService.GetVideosLikes([]int64{video.ID})
			// log.Println("获取点赞总数运行时间:", time.Since(t1))
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
			// t1 := time.Now()
			is_likes, err := videoService.AreVideosLikedByUser(user_id, []int64{video.ID})
			// log.Println("判断是否点赞运行时间:", time.Since(t1))
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
			// t1 := time.Now()
			temp_response_video.Id = video.ID
			temp_response_video.Play_url = video.PlayURL
			temp_response_video.Cover_url = video.CoverURL
			temp_response_video.Title = video.Title
			// log.Println("其他运行时间:", time.Since(t1))
			wait_group.Done()
		}(video, &temp_response_video)
		wait_group.Wait()
		// log.Println("运行时间:", time.Since(t_total))
		response_video_list[index] = temp_response_video
	}
	return response_video_list, nil
}

// 提供给外部的接口
func (videoService *VideoServiceImp) GetVideoListByAuthorID(authorId int64) ([]*model.Video, error) {
	dao_video_list, err := rabbitmq.DAOGetVideoListByAuthorID(authorId)
	if err != nil {
		return nil, err
	} else {
		return dao_video_list, nil
	}
}

// 提供给外部的接口
func (videoService *VideoServiceImp) GetVideoCountByAuthorID(authorId int64) (int64, error) {
	dao_video_list, err := rabbitmq.DAOGetVideoListByAuthorID(authorId)
	if err != nil {
		return 0, err
	} else {
		return int64(len(dao_video_list)), nil
	}
}

func (videoService *VideoServiceImp) PublishList(user_id int64) ([]ResponseVideo, error) {
	// 通过rabbitmq 获取数据库中的数据
	publishlistmq := rabbitmq.SimpleVideoPublishListMq
	dao_video_list, err := publishlistmq.PublishRequest("publishlist", time.Now().UnixMilli(), user_id)
	if err != nil {
		return nil, err
	}
	// dao.SetDefault(database.DB)
	// dao_video_list, err := videoService.GetVideoListByAuthorID(user_id)
	if err != nil {
		return nil, err
	}
	response_video_list, err := makeResponseVideo(dao_video_list, videoService, user_id)
	if err != nil {
		return nil, err
	}
	return response_video_list, err

}

// 上传视频接口新的
func (videoService *VideoServiceImp) Action(title string, userID int64, videoname string, file multipart.File) error {
	err := UploadVideoToOSS(videoname, file)
	if err != nil {
		log.Println("Upload Video ERROR : ", err)
		return err
	}
	err = InsertVideo(videoname, userID, title)
	if err != nil {
		log.Println("Insert Video ERROR : ", err)
		return err
	}
	return nil
}

// 将视频信息加入到数据库中
func InsertVideo(videoname string, userID int64, title string) error {
	var video model.Video
	playurl := oss1.URLPre + videoname
	video.Title = title
	video.AuthorID = userID
	video.PlayURL = playurl
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	video.CoverURL = playurl + oss1.CoverURL_SUFFIX
	Video := dao.Video
	err := Video.Create(&video)
	return err
}

// 视频上传到oss
func UploadVideoToOSS(videoname string, file multipart.File) error {
	log.Println("即将上传视频", videoname)
	err := oss1.Bucket.PutObject(videoname, file)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
	}
	return nil
}
func getbucketfile() {
	bucket := oss1.Bucket
	continueToken := ""
	for {
		lsRes, err := bucket.ListObjectsV2(oss.ContinuationToken(continueToken))
		if err != nil {
			log.Println("err", err.Error())
		}
		// 打印列举结果。默认情况下，一次返回100条记录。
		for _, object := range lsRes.Objects {
			fmt.Println(object.Key, object.Type, object.Size, object.ETag, object.LastModified, object.StorageClass)
		}
		if lsRes.IsTruncated {
			continueToken = lsRes.NextContinuationToken
		} else {
			break
		}
	}
}
