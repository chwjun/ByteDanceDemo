package rabbitmq

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var callback_chan = make(chan Response)

const size = 8

type VideoMQ struct {
	RabbitMQ
	QueueName string
}

type Request struct {
	UserID        int64
	Latest_time   int64
	CorrelationId string
}

type Response struct {
	CorrelationId     string
	ResponseVideoList []*model.Video
}

// 编码请求
func Encode_Request(request Request) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(request); err != nil {
		return nil, err
	}

	// 返回缓冲区中的字节切片
	return buf.Bytes(), nil
}

// 解码请求
func Decode_Request(data []byte) (Request, error) {
	// 创建一个字节缓冲区并将数据写入其中
	buf := bytes.NewBuffer(data)

	// 创建一个解码器，并从缓冲区中读取数据
	decoder := gob.NewDecoder(buf)

	// 创建一个变量来接收解码的结构体数组
	var request Request

	// 使用解码器将缓冲区中的数据解码为结构体数组
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}

	// 返回解码后的结构体数组
	return request, nil
}

func Encode_Response(response Response) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(response); err != nil {
		return nil, err
	}

	// 返回缓冲区中的字节切片
	return buf.Bytes(), nil
}

func Decoed_Response(data []byte) (Response, error) {
	// 创建一个字节缓冲区并将数据写入其中
	buf := bytes.NewBuffer(data)

	// 创建一个解码器，并从缓冲区中读取数据
	decoder := gob.NewDecoder(buf)

	// 创建一个变量来接收解码的结构体数组
	var response Response

	// 使用解码器将缓冲区中的数据解码为结构体数组
	if err := decoder.Decode(&response); err != nil {
		return response, err
	}

	// 返回解码后的结构体数组
	return response, nil
}

var SimpleVideoFeedMq *VideoMQ
var SimpleVideoPublishListMq *VideoMQ

func InitVideoRabbitMQ() {
	// 初始化MQ对象，订阅通道
	SimpleVideoFeedMq = NewSimpleVideoRabbitMQ("feed", "")
	SimpleVideoPublishListMq = NewSimpleVideoRabbitMQ("publishlist", "")
	log.Println("videoMQ")
	// 开启 go routine 启动消费者
	go SimpleVideoFeedMq.Consumer()
	go SimpleVideoPublishListMq.Consumer()
}

// NewSimpleVideoRabbitMQ 新建简单模式的消息队列（生产者，消息队列，一个消费者）
func NewSimpleVideoRabbitMQ(queueName string, replyto string) *VideoMQ {
	return newVideoRabbitMQ(queueName)
}

// 新建队列，初始化队列
func newVideoRabbitMQ(queuename string) *VideoMQ {
	videoMQ := &VideoMQ{
		RabbitMQ:  *BaseRmq,
		QueueName: queuename,
	}
	// log.Println(queuename, "there")
	queue, err := videoMQ.video_channel.QueueDeclare(
		queuename,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())

	}
	log.Println("queuename ", queue.Name)
	log.Println(queuename, "here")
	return videoMQ
}

// 发送消息
func (r *VideoMQ) PublishRequest(queuename string, latest_time int64, user_id int64) ([]*model.Video, error) {

	// var wg sync.WaitGroup
	corrID := uuid.New().String()
	// 声明回调队列
	callback_queue, err := r.video_channel.QueueDeclare(
		corrID,
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		return nil, err
	}
	callback_msgs, err := r.video_channel.Consume(
		callback_queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		return nil, err
	}
	// 发布消息
	var request Request
	request.CorrelationId = corrID
	request.Latest_time = latest_time
	request.UserID = user_id
	requestbuf, err := Encode_Request(request)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		return nil, err
	}
	log.Println("发送的消息", request.Latest_time, request.UserID)
	// 发送消息
	err = r.video_channel.Publish(
		"",
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ReplyTo:       corrID,
			CorrelationId: corrID,
			ContentType:   "application/octet-stream",
			Body:          requestbuf,
		},
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		return nil, err
	}
	// 等待消息确认
	log.Println("------------------------正在等待回调消息----------------------------")
	for d := range callback_msgs {
		log.Println("--------------获取回调消息----------------")
		var response Response
		response, err = Decoed_Response(d.Body)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("[%s : %d] %s \n", file, line, err.Error())
			return nil, err
		}
		log.Println("回调消息", response.CorrelationId, " ", len(response.ResponseVideoList))
		if response.CorrelationId == corrID {
			d.Ack(false)
			log.Println("成功获取消息的回调")
			// 删除队列
			_, err := r.video_channel.QueueDelete(
				corrID, // 队列名称
				false,  // 如果为 true，则只能删除被未被使用的队列，否则会返回错误
				false,  // 如果为 true，则只删除当前绑定的队列
				false,  // 如果为 true，则在等待服务器响应时删除队列
			)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			return response.ResponseVideoList, nil
		} else {
			d.Nack(false, true)
		}
	}
	// wg.Wait()
	return nil, err
}

// 消费者
func (r *VideoMQ) Consumer() {
	msgs, err := r.video_channel.Consume(
		r.QueueName,
		"",
		false, // 手动确认
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
	}
	for d := range msgs {
		log.Println("--------------------消费数据中----------------------")
		request, err := Decode_Request(d.Body)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("[%s : %d] %s \n", file, line, err.Error())
			continue
		}
		log.Println("获取到的数据", request.CorrelationId, " ", request.Latest_time, " ", request.UserID)
		var response Response
		// 如果是feed的数据
		if r.QueueName == "feed" {
			latest_time := time.UnixMilli(request.Latest_time)
			response_video_list, err := GetVideosByLatestTime(latest_time)
			if err != nil {
				_, file, line, _ := runtime.Caller(0)
				log.Printf("[%s : %d] %s \n", file, line, err.Error())
			}
			response.ResponseVideoList = response_video_list
			response.CorrelationId = request.CorrelationId
		}
		// 是publish list的请求
		if r.QueueName == "publishlist" {
			response_video_list, err := DAOGetVideoListByAuthorID(request.UserID)
			if err != nil {
				_, file, line, _ := runtime.Caller(0)
				log.Printf("[%s : %d] %s \n", file, line, err.Error())
			}
			response.ResponseVideoList = response_video_list
			response.CorrelationId = request.CorrelationId
		}
		responsebuf, err := Encode_Response(response)

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("[%s : %d] %s \n", file, line, err.Error())
		}

		err = r.video_channel.Publish(
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/octet-stream",
				CorrelationId: d.CorrelationId,
				Body:          responsebuf,
			},
		)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("[%s : %d] %s \n", file, line, err.Error())
		}
		d.Ack(false)
		log.Println("-------------------------数据消费完毕-------------------------")
	}
}

// 这个是video专用的通过时间获取videolist
func GetVideosByLatestTime(latest_time time.Time) ([]*model.Video, error) {
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
	log.Println("从数据库获取视频的数量", len(result))
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
