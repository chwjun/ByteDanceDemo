package rabbitmq

// "bytedancedemo/service"
// "bytedancedemo/service"

// "github.com/rs/xid"

import (
	"log"
	"runtime"

	amqp "github.com/rabbitmq/amqp091-go"
)

type VideoMQ struct {
	RabbitMQ
	// ch        *amqp.Channel
	QueueName string
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
func (r *VideoMQ) PublishRequest(queuename string) error {
	err := r.video_channel.Publish(
		"",
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(queuename),
		},
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
		return err
	}
	log.Println("消息发布")
	confirmations := r.video_channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	confirm := <-confirmations
	if confirm.Ack {
		log.Println("消息被确认")
		return nil
	} else {
		log.Println("消息被拒绝")
	}
	return nil
}

// 消费者
func (r *VideoMQ) Consumer() {
	q, err := r.video_channel.QueueDeclare(
		r.QueueName,
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
	msgs, err := r.video_channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s : %d] %s \n", file, line, err.Error())
	}
	for msg := range msgs {
		log.Printf("消费了msg : %s \n", msg.Body)
		msg.Ack(false)
	}
}
