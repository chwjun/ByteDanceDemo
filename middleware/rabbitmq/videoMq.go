package rabbitmq

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type VideoMQ struct {
	RabbitMQ
	QueueName string
	Exchange  string
	Key       string
	CorrID    string
	ReplyName string
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
	return videoMQ
}

// PublishSimpleVideo simple模式下视频请求生产者
func (r *VideoMQ) PublishSimpleVideo(message string, c *gin.Context) error {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
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
		log.Println(err)
		return err
	}
	//调用channel 发送消息到队列中
	err = r.channel.PublishWithContext(
		c,
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: r.CorrID,
			ReplyTo:       r.ReplyName,
			Body:          []byte(message),
		})
	if err != nil {
		return err
	}
	// 处理回调,判断消息是否处理
	content,err := r.channel.Consume(
		r.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for msg := range content {
		if msg.CorrelationId == r.CorrID && string(msg.Body) == "Finish" {
			// false表示仅确认当前消息
			msg.Ack(false)
			break
		}
	}
	return nil
}

// ConsumeSimpleVideo simple模式下消费者 video模块
func (r *VideoMQ) ConsumeSimpleVideo() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
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
		fmt.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
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
		fmt.Println(err)
	}

	forever := make(chan bool)
	//启用协程处理消息
	//go func() {
	//	for d := range msgs {
	//		//消息逻辑处理，可以自行设计逻辑
	//		log.Printf("Received a message: %s", d.Body)
	//	}
	//}()

	log.Println("q.Name", q.Name)
	switch q.Name {
	case "video_feed":
		go r.consumerVideoFeed(msgs)
	case "video_publish_list":
		go r.consumerVideoPublishList(msgs)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// consumerVideoFeed 添加Feed的消费，并且发送给生产者完成消息
func (r *VideoMQ) consumerVideoFeed(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		// 解析参数
		params := strings.Split(fmt.Sprintf("%s", msg.Body), "-")
		log.Println("添加视频消费者获得 params:", params)
		// 回调队列
		r.channel.Publish(
			r.Exchange,
			r.Key,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: r.CorrID,
				ReplyTo:       r.ReplyName,
				Body:          []byte("Finish"),
			},
		)
	}
}

// consumerVideoPublishList 添加publislist 的消费，并且发送给生产者完成消息
func (r *VideoMQ) consumerVideoPublishList(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		// 解析参数
		params := strings.Split(fmt.Sprintf("%s", msg.Body), "-")
		log.Println("添加视频消费者获得 params:", params)
		// 回调队列
		r.channel.Publish(
			r.Exchange,
			r.Key,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: r.CorrID,
				ReplyTo:       r.ReplyName,
				Body:          []byte("Finish"),
			},
		)
	}
}

// NewSimpleVideoRabbitMQ 新建简单模式的消息队列（生产者，消息队列，一个消费者）
func NewSimpleVideoRabbitMQ(queueName string, replyto string, corrid string) *VideoMQ {
	return newVideoRabbitMQ(queueName, "", "", replyto, corrid)
}
func InitVideoRabbitMQ() {
	SimpleVideoFeedMq = NewSimpleVideoRabbitMQ("feed", "feed", "feed")
	SimpleVideoPublishListMq = NewSimpleVideoRabbitMQ("publishlist", "publishlist", "publishlist")
	// 开启 go routine 启动消费者
	go SimpleVideoFeedMq.ConsumeSimpleVideo()
	go SimpleVideoPublishListMq.ConsumeSimpleVideo()
}
