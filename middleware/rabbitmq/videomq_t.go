package rabbitmq

import (
	"bytes"
	"encoding/gob"
	"log"
	"runtime"
)

var callback_chan = make(chan Response)

type User struct {
	Id              int64  `json:"id"`               // 主键
	Name            string `json:"name"`             // 用户名 用于登录 不可重复
	FollowCount     int64  `json:"follow_count"`     // 关注数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝数
	IsFollow        bool   `json:"is_follow"`        // 是否关注
	Avatar          string `json:"avatar"`           // 头像
	BackgroundImage string `json:"background_image"` // 背景图
	Signature       string `json:"signature"`        // 个性签名
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数
	WorkCount       int64  `json:"work_count"`       // 作品数
	FavoriteCount   int64  `json:"favorite_count"`   // 点赞数
}
type ResponseVideo struct {
	Id int64 `json:"id,omitempty"`
	// 作者信息
	Author         User   `json:"author"`
	Play_url       string `json:"play_url" json:"play_url,omitempty"`
	Cover_url      string `json:"cover_url,omitempty"`
	Favorite_count int64  `json:"favorite_count,omitempty"`
	Comment_count  int64  `json:"comment_count,omitempty"`
	Is_favorite    bool   `json:"is_favorite,omitempty"`
	Title          string `json:"title,omitempt"`
}

type VideoMQ_t struct {
	RabbitMQ
	// ch        *amqp.Channel
	QueueName string
}

type Request struct {
	UserID        int64
	Latest_time   int64
	CorrelationId string
}

type Response struct {
	ResponseVideoList []ResponseVideo
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

var SimpleVideoFeedMq_t *VideoMQ
var SimpleVideoPublishListMq_t *VideoMQ

func InitVideoRabbitMQ_t() {
	// 初始化MQ对象，订阅通道
	SimpleVideoFeedMq = NewSimpleVideoRabbitMQ("feed", "")
	SimpleVideoPublishListMq = NewSimpleVideoRabbitMQ("publishlist", "")
	log.Println("videoMQ")
	// 开启 go routine 启动消费者
	go SimpleVideoFeedMq.Consumer()
	go SimpleVideoPublishListMq.Consumer()
}

// NewSimpleVideoRabbitMQ 新建简单模式的消息队列（生产者，消息队列，一个消费者）
func NewSimpleVideoRabbitMQ_t(queueName string, replyto string) *VideoMQ {
	return newVideoRabbitMQ_t(queueName)
}

// 新建队列，初始化队列
func newVideoRabbitMQ_t(queuename string) *VideoMQ {
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
func (r *VideoMQ) PublishRequest_t(queuename string) error {

}

// 消费者
func (r *VideoMQ) Consumer_t() {

}
