package service

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gookit/slog"

	"github.com/RaymondCode/simple-demo/model"

	"github.com/RaymondCode/simple-demo/dao"

	"gorm.io/gorm"

	"github.com/RaymondCode/simple-demo/util"
)

type Task struct {
	UserID     uint
	VideoID    uint
	Action     int32 // 1 for like, 2 for unlike
	RetryCount int
	ResultChan chan<- Result // 用于返回结果的通道
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
	slog.Debug("NewTaskQueue")
	return &TaskQueue{
		tasks: make([]*Task, 0),
	}
}

func (q *TaskQueue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	length := len(q.tasks)
	slog.Debug("TaskQueue length:", length)
	return length
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
func (q *TaskQueue) PushTask(task *Task, dispatchSignal chan bool) {
	heap.Push(q, task) // 使用标准的堆操作
	slog.Debug("Pushed task to queue: %+v", task)
	slog.Debug("TaskQueue length after push:", len(q.tasks)) // 添加此行
	dispatchSignal <- true                                   // 发送信号
	slog.Debug("Dispatch signal sent")
}

func (q *TaskQueue) PopTask() *Task {
	task := heap.Pop(q).(*Task) // 使用标准的堆操作
	slog.Debug("Popped task from queue: %+v", task)
	return task
}

func (q *TaskQueue) Push(x interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	item, ok := x.(*Task)
	if !ok {
		slog.Fatal("Push: Expected *Task, got something else")
		return
	}
	q.tasks = append(q.tasks, item)
}

func (q *TaskQueue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	n := len(q.tasks)
	item := q.tasks[n-1]
	q.tasks = q.tasks[0 : n-1]
	return item
}

func startWorkers(workerCount int, tasks *TaskQueue, results chan<- Result, quit chan bool, dispatchSignal chan bool) {
	var wg sync.WaitGroup
	taskCh := make(chan Task)

	slog.Debug("Starting task dispatcher")

	go func() {
		for {
			select {
			case _, ok := <-dispatchSignal:
				if !ok {
					// dispatchSignal 被关闭，退出循环
					return
				}
				slog.Debug("Received dispatch signal")
				if tasks.Len() > 0 {
					task := tasks.PopTask()
					slog.Debug("Dispatching task:", task)
					taskCh <- *task
				} else {
					slog.Debug("No tasks in queue")
				}
			default:
				time.Sleep(time.Millisecond * 100) // sleep for a while before next check
			}
		}
	}()
	slog.Debug("Starting workers")

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		slog.Debug("Starting worker", i)
		go worker(taskCh, results, quit, tasks, &wg, dispatchSignal)
	}

	wg.Wait() // 等待所有工作协程完成

	slog.Debug("All workers have finished")
}

func shouldRetry(task Task) bool {
	// 重试最多3次
	if task.RetryCount < 3 {
		return true
	}
	return false
}
func worker(tasks <-chan Task, results chan<- Result, quit <-chan bool, taskQueue *TaskQueue, wg *sync.WaitGroup, dispatchSignal chan bool) {
	defer wg.Done() // 在工作协程结束时调用
	for {

		select {
		case <-quit:
			slog.Debug("Worker stopped.")
			return
		case task := <-tasks:
			slog.Debug("Processing task: %+v", task)
			result := processTask(task)
			if result.Error != nil && shouldRetry(task) {
				task.RetryCount++                                        // 增加重试次数
				time.Sleep(time.Second * time.Duration(task.RetryCount)) // 增加延迟
				// 重新将任务推入堆中
				taskQueue.PushTask(&task, dispatchSignal)
			} else {
				slog.Debug("Sending result to results channel: %+v", result)
				task.ResultChan <- result
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
		err = util.GlobalRedisClient.LikeVideo(task.UserID, task.VideoID)
	case 2:
		err = util.GlobalRedisClient.UnlikeVideo(task.UserID, task.VideoID)
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

func SyncLikesToDatabase() error {
	videoIDs, err := getAllVideoIDs() // 获取所有视频ID
	if err != nil {
		return err
	}

	for _, videoID := range videoIDs {
		likes, err := util.GlobalRedisClient.GetLikes(videoID)
		if err != nil {
			return err
		}

		for userIDString, likedString := range likes {
			userID, err := strconv.Atoi(userIDString)
			if err != nil {
				return err // 或者可以记录错误并继续
			}
			liked, err := strconv.Atoi(likedString)
			if err != nil {
				return err // 或者可以记录错误并继续
			}

			// 开始一个新的事务
			tx := dao.DB.Begin()
			if tx.Error != nil {
				return tx.Error
			}

			first, err := dao.Like.Where(dao.Like.UserID.Eq(uint(userID)), dao.Like.VideoID.Eq(uint(videoID))).First()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 如果记录未找到，则创建一个新的喜欢记录
					like := model.Like{
						UserID:  uint(userID),
						VideoID: uint(videoID),
						Liked:   int(uint(liked)),
					}
					if err := tx.Create(&like).Error; err != nil {
						tx.Rollback() // 回滚事务
						return err
					}
				} else {
					// 如果发生其他错误，则回滚事务并返回该错误
					tx.Rollback()
					return err
				}
			} else {
				// 如果记录存在，则更新点赞状态
				first.Liked = int(uint(liked))
				if err := tx.Save(&first).Error; err != nil {
					tx.Rollback() // 回滚事务
					return err
				}
			}

			// 提交事务
			if err := tx.Commit().Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func getAllVideoIDs() ([]uint, error) {
	videos := dao.Video // 假设您有一个名为dao的包，其中有一个Video的DAO对象

	var videoIDs []uint
	ids, err := videos.Select(videos.ID).Find()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		videoIDs = append(videoIDs, id.ID) // 假设ID是Video中的一个字段
	}

	return videoIDs, nil
}

func StartSyncTask(db *gorm.DB, syncInterval time.Duration) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := SyncLikesToDatabase(); err != nil {
				log.Printf("Failed to sync likes to database: %v", err)
			}
		}
	}
}

var (
	taskQueue      *TaskQueue
	dispatchSignal chan bool
	results        chan Result
	quit           chan bool
)

func (s *FavoriteServiceImpl) StartFavoriteAction() {
	// 创建全局变量
	taskQueue = NewTaskQueue()
	dispatchSignal = make(chan bool, 10)
	results = make(chan Result, 10)
	quit = make(chan bool)

	// 初始化堆
	heap.Init(taskQueue)

	// 启动工人 Usmups.
	go startWorkers(10, taskQueue, results, quit, dispatchSignal)

}
func (s *FavoriteServiceImpl) FavoriteAction(userId int64, videoID int64, actionType int32) (FavoriteActionResponse, error) {
	slog.Debug("Starting FavoriteAction for videoID:", videoID, "actionType:", actionType)
	//dispatchSignal := make(chan bool, 10)
	//taskQueue := NewTaskQueue()
	//slog.Debug("TaskQueue created:", taskQueue) // 添加此行来查看任务队列的详细信息
	//// 初始化堆
	//heap.Init(taskQueue)
	//results := make(chan Result, 10)

	task := &Task{
		UserID:     uint(userId),
		VideoID:    uint(videoID),
		Action:     actionType,
		ResultChan: results, // 将结果通道添加到任务中
	}
	slog.Debug("Task created:", task) // 添加此行来查看任务的详细信息
	taskQueue.PushTask(task, dispatchSignal)
	slog.Debug("Task pushed to queue")

	//quit := make(chan bool)

	// 启动工人
	//go startWorkers(5, taskQueue, results, quit, s.util.GlobalRedisClient, dispatchSignal)
	//
	//slog.Debug("Workers started")
	//
	//slog.Debug("Task pushed to the queue:", task)

	// 等待结果
	slog.Debug("Waiting for results")
	result := <-results

	slog.Debug("Results received:")

	// 关闭任务通道
	//if taskQueue.Len() == 0 {
	//	close(quit)
	//	close(dispatchSignal)
	//}
	slog.Debug("FavoriteAction completed for videoID:", videoID)
	slog.Infof("Number of goroutines: %d\n", runtime.NumGoroutine())
	return FavoriteActionResponse{
		StatusCode: int32(result.StatusCode),
		StatusMsg:  result.StatusMsg,
	}, result.Error
}
