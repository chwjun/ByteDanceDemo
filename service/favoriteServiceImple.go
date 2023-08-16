package service

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"

	"github.com/RaymondCode/simple-demo/dao"
)

type FavoriteServiceImpl struct {
}
type Task struct {
	UserID     uint
	VideoID    uint
	Action     int32 // 1 for like, 2 for unlike
	RetryCount int
}

type Result struct {
	Success    bool
	Error      error
	StatusCode int
	StatusMsg  string
}
type TaskQueue struct {
	tasks []*Task
	mutex sync.Mutex
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks: make([]*Task, 0),
	}
}

func (q *TaskQueue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.tasks)
}

func (q *TaskQueue) Less(i, j int) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.tasks[i].Action < q.tasks[j].Action
}

func (q *TaskQueue) Swap(i, j int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks[i], q.tasks[j] = q.tasks[j], q.tasks[i]
}
func (q *TaskQueue) PushTask(task *Task) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	heap.Push(q, task)
}

func (q *TaskQueue) PopTask() *Task {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return heap.Pop(q).(*Task)
}

func (q *TaskQueue) Push(x interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	item := x.(*Task)
	q.tasks = append(q.tasks, item)
}

func (q *TaskQueue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	old := q.tasks
	n := len(old)
	item := old[n-1]
	q.tasks = old[0 : n-1]
	return item
}

// 实现heap.Interface
func startWorkers(workerCount int, tasks *TaskQueue, results chan<- Result, quit chan bool) {
	var wg sync.WaitGroup
	taskCh := make(chan Task)
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				if tasks.Len() > 0 { // 检查是否还有任务
					task := tasks.PopTask() // 使用新的PopTask方法
					taskCh <- *task
				}
			}
		}
	}()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(taskCh, results, quit, tasks, &wg)
	}

	wg.Wait() // 等待所有工作协程完成
}

func shouldRetry(task Task) bool {
	// 重试最多3次
	if task.RetryCount < 3 {
		return true
	}
	return false
}
func worker(tasks <-chan Task, results chan<- Result, quit <-chan bool, taskQueue *TaskQueue, wg *sync.WaitGroup) {
	defer wg.Done() // 在工作协程结束时调用
	for {

		select {
		case <-quit:
			return
		case task := <-tasks:
			result := processTask(task)
			if result.Error != nil && shouldRetry(task) {
				task.RetryCount++                                        // 增加重试次数
				time.Sleep(time.Second * time.Duration(task.RetryCount)) // 增加延迟
				// 重新将任务推入堆中
				taskQueue.PushTask(&task)
			} else {
				results <- result
			}
		}
	}
}

func processTask(task Task) Result {
	var err error
	statusCode := SuccessCode
	statusMsg := SuccessMessage

	switch task.Action {
	case 1:
		err = likeVideo(task.UserID, task.VideoID)
		if err != nil {
			statusCode = ErrorCode
			statusMsg = fmt.Sprintf("Failed to like video: %v", err)
		} else {
			statusCode = SuccessCode
			statusMsg = SuccessMessage
		}
	case 2:
		err = unlike(task.UserID, task.VideoID)
		if err != nil {
			statusCode = ErrorCode
			statusMsg = fmt.Sprintf("Failed to unlike video: %v", err)
		} else {
			statusCode = SuccessCode
			statusMsg = SuccessMessage
		}
	default:
		err = fmt.Errorf("invalid action_type: %v", task.Action)
		statusCode = ErrorCode
		statusMsg = err.Error()
	}

	return Result{
		Success:    err == nil,
		Error:      err,
		StatusCode: int(statusCode),
		StatusMsg:  statusMsg,
	}
}

func (s *FavoriteServiceImpl) FavoriteAction(videoID int64, actionType int32) (FavoriteActionResponse, error) {
	taskQueue := NewTaskQueue()
	heap.Init(taskQueue)

	results := make(chan Result, 10)
	quit := make(chan bool)

	// 启动工人
	startWorkers(5, taskQueue, results, quit)

	task := &Task{
		UserID:  1,
		VideoID: uint(videoID),
		Action:  actionType,
	}
	taskQueue.PushTask(task)

	// 等待结果
	result := <-results

	// 关闭任务通道
	close(quit)

	return FavoriteActionResponse{
		StatusCode: int32(result.StatusCode),
		StatusMsg:  result.StatusMsg,
	}, result.Error
}

func (s *FavoriteServiceImpl) FavoriteList(userID int64) (FavoriteListResponse, error) {
	// 通过userId获取用户点赞的视频列表
	videoList, err := s.GetFavoriteVideoInfoByUserID(userID)
	if err != nil {
		errorCode := ErrorCode
		errorMessage := "获取视频失败: " + err.Error()
		return FavoriteListResponse{
			StatusCode: errorCode,
			StatusMsg:  errorMessage,
		}, nil
	}

	successCode := SuccessCode
	successMessage := SuccessMessage
	response := FavoriteListResponse{
		StatusCode: successCode,
		StatusMsg:  successMessage,
		VideoList:  videoList,
	}

	return response, nil
}

func (s *FavoriteServiceImpl) GetFavoriteVideoInfoByUserID(userID int64) ([]*Video, error) {
	videoIDs, err := GetLikedVideoIDs(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("获取点赞视频ID失败: %v", err)
	}

	var videos []*Video
	for _, videoID := range videoIDs {
		// 使用特定的查询构造方式获取视频详情
		authorID, title, playURL, coverURL, err := GetVideoDetailsByID(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取视频详情失败: %v", err)
		}

		commentCount, err := GetCommentCount(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取评论总数失败: %v", err)
		}

		likeCount, err := GetLikeCount(videoID)
		if err != nil {
			return nil, fmt.Errorf("获取点赞总数失败: %v", err)
		}

		isFavorite, err := IsVideoLikedByUser(uint(userID), videoID)
		if err != nil {
			return nil, fmt.Errorf("判断用户是否点赞了视频失败: %v", err)
		}

		requestingUserID := int64(userID)
		author, err := s.GetUserInfoByID(&requestingUserID, int64(authorID))
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %v", err)
		}

		video := &Video{
			ID:            int64(videoID),
			Author:        *author,
			PlayURL:       playURL,
			CoverURL:      coverURL,
			FavoriteCount: likeCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         title,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func likeVideo(userID uint, videoID uint) error {
	// 开始一个新的事务
	tx := dao.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 使用特定的查询构造方式
	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
	if err != nil {
		// 如果记录未找到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建一个新的喜欢记录
			like := model.Like{
				UserID:  userID,
				VideoID: videoID,
				Liked:   1,
			}
			// 将新记录保存到数据库
			if err := tx.Create(&like).Error; err != nil {
				tx.Rollback() // 回滚事务
				return err
			}
		} else {
			// 如果发生其他错误，则回滚事务并返回该错误
			tx.Rollback()
			return err
		}
	}

	// 假设 first 是一个 *model.Like 类型
	if first.Liked == 1 {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("user has already liked this video")
	}

	// 将喜欢的状态设置为1
	first.Liked = 1
	// 保存记录
	if err := tx.Save(&first).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
func unlike(userID uint, videoID uint) error {
	// 开始一个新的事务
	tx := dao.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 使用特定的查询构造方式
	first, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.VideoID.Eq(videoID)).First()
	if err != nil {
		// 如果记录未找到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("No like found for this user and video")
		}
		// 如果发生其他错误，则回滚事务并返回该错误
		tx.Rollback()
		return err
	}

	// 假设 first 是一个 *model.Like 类型
	if first.Liked == 0 {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("User has already unliked this video")
	}

	// 将喜欢的状态设置为0
	first.Liked = 0
	// 保存记录
	if err := tx.Save(&first).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
func (s *FavoriteServiceImpl) GetUserInfoByID(requestingUserID *int64, userID int64) (*User, error) {
	// 使用特定的查询构造方式获取用户详情
	name, avatar, backgroundImage, signature, err := GetUserDetailsByID(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取关注总数
	followCount, err := GetUserFollowCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取粉丝总数
	followerCount, err := GetUserFollowerCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 检查是否已关注
	isFollow, err := IsUserFollowingAnotherUser(uint(*requestingUserID), uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取获赞总数
	totalFavorited, err := GetUserTotalReceivedLikes(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取作品数量
	workCount, err := GetUserWorkCount(uint(userID))
	if err != nil {
		return nil, err
	}

	// 获取点赞数量
	favoriteCount, err := GetUserTotalReceivedLikes(uint(userID))
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:              userID,
		Name:            name,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        isFollow,
		Avatar:          avatar,
		BackgroundImage: backgroundImage,
		Signature:       signature,
		TotalFavorited:  totalFavorited,
		WorkCount:       workCount,
		FavoriteCount:   favoriteCount,
	}

	return user, nil
}
func GetUserDetailsByID(userID uint) (string, string, string, string, error) {

	first, err := dao.User.Select(dao.User.Name, dao.User.Avatar, dao.User.BackgroundImage, dao.User.Signature).Where(dao.User.ID.Eq(userID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", "", "", fmt.Errorf("No user found for this user ID")
		}
	}
	return first.Name, first.Avatar, first.BackgroundImage, first.Signature, nil
}

func GetLikedVideoIDs(userID uint) ([]uint, error) {

	likes, err := dao.Like.Where(dao.Like.UserID.Eq(userID), dao.Like.Liked.Eq(1)).Order(dao.Like.CreatedAt.Abs()).Find()

	if err != nil {
		return nil, err
	}

	var videoIDs []uint
	for _, like := range likes {
		videoIDs = append(videoIDs, like.VideoID) // 假设VideoID是model.Like中的一个字段
	}

	return videoIDs, nil
}
func GetVideoDetailsByID(videoID uint) (uint, string, string, string, error) {

	First, err := dao.Video.Select(dao.Video.AuthorID, dao.Video.Title, dao.Video.PlayURL, dao.Video.CoverURL).Where(dao.Video.ID.Eq(videoID)).First()
	if err != nil {
		return 0, "", "", "", fmt.Errorf("找不到视频ID %d: %v", videoID, err)
	}
	return First.AuthorID, First.Title, First.PlayURL, First.CoverURL, nil

}
func GetCommentCount(videoID uint) (int64, error) {

	var count int64
	count, err := dao.Comment.Where(dao.Comment.VideoID.Eq(videoID), dao.Comment.ActionType.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func GetLikeCount(videoID uint) (int64, error) {

	count, err := dao.Like.Where(dao.Like.VideoID.Eq(videoID), dao.Like.Liked.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func GetUserFollowCount(userID uint) (int64, error) {

	count, err := dao.Relation.Where(dao.Relation.UserID.Eq(userID), dao.Relation.Followed.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func GetUserFollowerCount(userID uint) (int64, error) {

	count, err := dao.Relation.Where(dao.Relation.FollowingID.Eq(userID), dao.Relation.Followed.Eq(1)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
func IsUserFollowingAnotherUser(userID, followingID uint) (bool, error) {
	relation, err := dao.Relation.Where(dao.Relation.UserID.Eq(userID), dao.Relation.FollowingID.Eq(followingID), dao.Relation.Followed.Eq(1)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // 没有找到关系记录，表示用户没有关注另一个用户
		}
		return false, err // 其他错误
	}

	return relation.Followed == 1, nil // 根据Followed字段返回是否关注
}

func GetUserTotalReceivedLikes(userID uint) (int64, error) {
	likes := dao.Like
	videos := dao.Video

	count, err := likes.Join(videos, videos.ID.EqCol(likes.VideoID)).Where(videos.AuthorID.Eq(userID), likes.Liked.Eq(1)).Count()

	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserWorkCount(userID uint) (int64, error) {

	workCount, err := dao.Video.Where(dao.Video.AuthorID.Eq(userID)).Count()

	if err != nil {
		return 0, err
	}
	return workCount, nil
}

func IsVideoLikedByUser(userID uint, videoID uint) (bool, error) {
	likes := dao.Like

	count, err := likes.Where(likes.UserID.Eq(userID), likes.VideoID.Eq(videoID), likes.Liked.Eq(1)).Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
