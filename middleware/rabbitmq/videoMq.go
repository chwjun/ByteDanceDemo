package rabbitmq

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	retry = 10
)

type VideoMQ struct {
	RabbitMQ
	QueueName string
	Exchange  string
	Key       string
	CorrID    string
	ReplyName string
	// CallBackExchange string
	Call_Back_Queue amqp.Queue
}

var SimpleVideoFeedMq *VideoMQ
var SimpleVideoPublishListMq *VideoMQ

// 新建“视频”消息队列
func newVideoRabbitMQ(queueName string, exchangeName string, key string, replyto string, corrid string) *VideoMQ {
	videoMQ := &VideoMQ{
		RabbitMQ:  *BaseRmq,
		QueueName: queueName,
		Exchange:  exchangeName,
		Key:       key,
		ReplyName: replyto,
		CorrID:    corrid,
	}
	var call_back_queue amqp.Queue
	var err error
	for i := 0; i < retry; i++ {
		// 声明回调队列
		call_back_queue, err = videoMQ.channel.QueueDeclare(
			replyto,
			//是否持久化
			false,
			//是否自动删除
			false,
			//是否具有排他性
			false,
			//是否阻塞处理
			false,
			//额外的属性
			nil,
		)
		if err != nil {
			log.Println("创建回调队列重试次数 ： ", i, queueName)
		} else {
			break
		}
	}
	if err != nil {
		log.Println("无法创建回调队列", err.Error(), queueName)
	}
	videoMQ.Call_Back_Queue = call_back_queue
	return videoMQ
}

// PublishSimpleVideo simple模式下视频请求生产者
func (r *VideoMQ) PublishSimpleVideo(message string, c *gin.Context) error {
	fmt.Println("video生产" + r.QueueName)
	//1.申请回调队列，如果队列不存在会自动创建，存在则跳过创建
	call_back_queue := r.Call_Back_Queue

	// 消费回调队列中的消息
	content, err := r.channel.Consume(
		call_back_queue.Name, // 回调队列名称
		"",                   // 消费者名称，为空时会随机生成一个名称
		false,                // 自动确认
		false,                // 独占
		false,                // 非阻塞
		false,                // 其他属性
		nil,                  // 额外属性
	)
	if err != nil {
		log.Println("无法消费回调队列中的消息", err)
	}
	// 发送请求消息到请求队列
	err = r.channel.Publish(
		"",
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			ReplyTo:     call_back_queue.Name,
		},
	)
	if err != nil {
		log.Println("无法发送请求消息:", err)
	}
	log.Println("发送请求")
	for msg := range content {
		fmt.Println("msg来了")
		if msg.CorrelationId == r.CorrID && string(msg.Body) == "Finish" {
			// false表示仅确认当前消息
			msg.Ack(false)
			break
		}
	}
	fmt.Println("回调处理完毕")
	return nil
}

// ConsumeSimpleVideo simple模式下消费者 video模块
func (r *VideoMQ) ConsumeSimpleVideo() {
	call_back_queue := r.Call_Back_Queue
	//接收消息
	msgs, err := r.channel.Consume(
		r.QueueName, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Connection中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println("接收消息失败")
	}

	forever := make(chan bool)
	//启用协程处理消息
	//go func() {
	//	for d := range msgs {
	//		//消息逻辑处理，可以自行设计逻辑
	//		log.Printf("Received a message: %s", d.Body)
	//	}
	//}()

	log.Println("call back queue.Name", call_back_queue.Name)
	switch call_back_queue.Name {
	case "feedreply":
		go r.consumerVideoFeed(msgs)
	case "publishlistreply":
		go r.consumerVideoPublishList(msgs)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// consumerVideoFeed 添加Feed的消费，并且发送给生产者完成消息
func (r *VideoMQ) consumerVideoFeed(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		fmt.Println("开始消费")
		// 解析参数
		params := strings.Split(fmt.Sprintf("%s", msg.Body), "-")
		log.Println("添加视频消费者获得 params:", params)
		log.Println("MQ参数：queue name : ", r.QueueName, " CorrID : ", r.CorrID)
		// 回调队列
		err := r.channel.Publish(
			"",
			msg.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: r.CorrID,
				// ReplyTo:       msg.ReplyTo,
				Body: []byte("Finish"),
			},
		)
		if err != nil {
			log.Println("消费后回调出错：", err)
		}
		fmt.Println("回调发送完毕")
	}
}

// consumerVideoPublishList 添加publislist 的消费，并且发送给生产者完成消息
func (r *VideoMQ) consumerVideoPublishList(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		fmt.Println("开始消费")
		// 解析参数
		params := strings.Split(fmt.Sprintf("%s", msg.Body), "-")
		log.Println("添加视频消费者获得 params:", params)
		log.Println("MQ参数：queue name : ", r.QueueName, " CorrID : ", r.CorrID)
		// 回调队列
		err := r.channel.Publish(
			"",
			msg.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: r.CorrID,
				// ReplyTo:       msg.ReplyTo,
				Body: []byte("Finish"),
			},
		)
		if err != nil {
			log.Println("消费后回调出错：", err)
		}
		fmt.Println("回调发送完毕")
	}
}

// NewSimpleVideoRabbitMQ 新建简单模式的消息队列（生产者，消息队列，一个消费者）
func NewSimpleVideoRabbitMQ(queueName string, replyto string, corrid string) *VideoMQ {
	return newVideoRabbitMQ(queueName, "", "", replyto, corrid)
}
func InitVideoRabbitMQ() {
	// 初始化MQ对象，订阅通道
	SimpleVideoFeedMq = NewSimpleVideoRabbitMQ("feed", "feedreply", "feed")
	SimpleVideoPublishListMq = NewSimpleVideoRabbitMQ("publishlist", "publishlistreply", "publishlist")
	// 开启 go routine 启动消费者
	go SimpleVideoFeedMq.ConsumeSimpleVideo()
	go SimpleVideoPublishListMq.ConsumeSimpleVideo()
}
