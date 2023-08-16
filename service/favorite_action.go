package service

import (
	"container/heap"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/RaymondCode/simple-demo/util"
)

type FavoriteAction struct {
	redisClient *util.RedisClient
}

func NewFavoriteService(redisAddr string, redisPassword string, redisDB int) *FavoriteServiceImpl {
	return &FavoriteServiceImpl{
		redisClient: util.NewRedisClient(redisAddr, redisPassword, redisDB),
	}
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
func startWorkers(workerCount int, tasks *TaskQueue, results chan<- Result, quit chan bool, redisClient *util.RedisClient) {
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
		go worker(taskCh, results, quit, tasks, &wg, redisClient)
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
func worker(tasks <-chan Task, results chan<- Result, quit <-chan bool, taskQueue *TaskQueue, wg *sync.WaitGroup, redisClient *util.RedisClient) {
	defer wg.Done() // 在工作协程结束时调用
	for {

		select {
		case <-quit:
			return
		case task := <-tasks:
			result := processTask(task, redisClient)
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

func processTask(task Task, redisClient *util.RedisClient) Result {
	var err error
	statusCode := SuccessCode
	statusMsg := SuccessMessage

	switch task.Action {
	case 1:
		err = redisClient.IncrementLikes(task.VideoID)
	case 2:
		err = redisClient.DecrementLikes(task.VideoID)
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

// 其他代码不变

func startSyncTask(redisClient *util.RedisClient, syncInterval time.Duration) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	syncFunc := func(videoID uint, likes int64) error {
		// 这里填写同步到数据库的逻辑
		// 例如，更新数据库中对应视频的点赞数
		return nil
	}

	for {
		select {
		case <-ticker.C:
			if err := redisClient.SyncLikesToDatabase(syncFunc); err != nil {
				log.Printf("Failed to sync likes to database: %v", err)
			}
		}
	}
}
func getAllVideoIDs() []uint {
	// 这里您可以从数据库或其他存储中获取所有的视频ID
	// 返回一个uint类型的切片
	return []uint{} // 示例返回空切片
}

func (s *FavoriteServiceImpl) FavoriteAction(videoID int64, actionType int32) (FavoriteActionResponse, error) {
	taskQueue := NewTaskQueue()
	heap.Init(taskQueue)

	results := make(chan Result, 10)
	quit := make(chan bool)

	// 启动工人
	startWorkers(5, taskQueue, results, quit, s.redisClient)

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
