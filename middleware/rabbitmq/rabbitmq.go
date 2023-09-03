package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

var BaseRmq *RabbitMQ
var MQURL string

// RabbitMQ 参考: https://www.rabbitmq.com/tutorials/tutorial-one-go.html
type RabbitMQ struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	video_conn    *amqp.Connection
	video_channel *amqp.Channel
	// 连接信息
	MqUrl string
}

func InitRabbitMQ() {
	MQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/",
		viper.GetString("settings.rabbitMQ.username"),
		viper.GetString("settings.rabbitMQ.password"),
		viper.GetString("settings.rabbitMQ.host"),
		viper.GetString("settings.rabbitMQ.port"),
	)
	BaseRmq = &RabbitMQ{
		MqUrl: MQURL,
	}
	conn, err := amqp.Dial(BaseRmq.MqUrl)
	if err != nil {
		BaseRmq.failOnError(err, "Failed to connect to RabbitMQ")
		log.Println("无法连接到rabbitmq")
	}
	conn1, err := amqp.Dial(BaseRmq.MqUrl)
	if err != nil {
		BaseRmq.failOnError(err, "Failed to connect to RabbitMQ")
		log.Println("无法连接到rabbitmq")
	}
	BaseRmq.conn = conn
	BaseRmq.video_conn = conn1
	BaseRmq.channel, err = conn.Channel()
	BaseRmq.failOnError(err, "Failed to get channel")
	BaseRmq.video_channel, err = conn1.Channel()
	if err != nil {
		BaseRmq.failOnError(err, "Failed to set to confirm mode")
		log.Println("无法设置confirm模式")
	}
	// err = BaseRmq.video_channel.Confirm(false)
	// if err != nil {
	// 	BaseRmq.failOnError(err, "Failed to set to confirm mode")
	// 	log.Println("无法设置confirm模式")
	// }
	log.Println("RabbitMQ")
}

func (r *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		log.Panicf("%s: %s", msg, err)
	}
}

func (r *RabbitMQ) destroy() {
	r.channel.Close()
	r.conn.Close()
}
