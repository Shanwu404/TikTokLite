package rabbitmq

import (
	"log"
	"strconv"
	"strings"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/streadway/amqp"
)

type RabbitMQLike struct {
	RabbitMQ
	Channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	// routing Key
	RoutingKey string
}

var RabbitMQLikeAdd *RabbitMQLike
var RabbitMQLikeDel *RabbitMQLike

// 简单模式 点赞mq实例
func NewLikeRabbitMQ(queueName string) *RabbitMQLike {
	rbq := &RabbitMQLike{
		RabbitMQ:   *MyRabbitMQ,
		QueueName:  queueName,
		Exchange:   "",
		RoutingKey: "",
	}
	var err error
	rbq.Channel, err = rbq.Conn.Channel()
	MyRabbitMQ.failOnErr(err, "Failed to get channel")
	return rbq
}

func InitLikeRabbitMQ() {
	RabbitMQLikeAdd = NewLikeRabbitMQ("Like Add")
	go RabbitMQLikeAdd.Consumer()
	log.Println("RabbitMQLikeAdd init successfully!")
	RabbitMQLikeDel = NewLikeRabbitMQ("Like Del")
	go RabbitMQLikeDel.Consumer()
	log.Println("RabbitMQLikeDel init successfully!")
}

func (likemq *RabbitMQLike) Producer(message string) {
	// 1. 申请消息队列
	_, err := likemq.Channel.QueueDeclare(
		likemq.QueueName,
		false, //是否持久化
		false, //是否为自动删除
		false, //是否具有排他性
		false, //是否阻塞
		nil,   //额外属性
	)
	if err != nil {
		log.Println(err)
	}

	//2. 发送消息到消息队列
	likemq.Channel.Publish(
		likemq.Exchange,
		likemq.QueueName,
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(message)},
	)
}

func (likemq *RabbitMQLike) Consumer() {
	// 1. 申请消息队列，如果队列不存在则自动创建，如果存在则跳过创建
	_, err := likemq.Channel.QueueDeclare(
		likemq.QueueName,
		false, //是否持久化
		false, //是否为自动删除
		false, //是否具有排他性
		false, //是否阻塞
		nil,   //额外属性
	)
	if err != nil {
		log.Println(err)
	}
	//2. 接收消息
	messages, err := likemq.Channel.Consume(
		likemq.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)
	if likemq.QueueName == "Like Add" {
		//处理点赞添加队列中的数据
		go consumerLikeAdd(messages)
	} else if likemq.QueueName == "Like Del" {
		//处理点赞删除队列中的数据
		go consumerLikeDel(messages)

	}
	log.Printf("[*] waiting for messages, to exit process CTRL+C")
	<-forever
}

func consumerLikeAdd(messages <-chan amqp.Delivery) {
	for message := range messages {
		//解析数据
		data := strings.Split(string(message.Body), ":")
		userId, _ := strconv.ParseInt(data[0], 10, 64)
		videoId, _ := strconv.ParseInt(data[1], 10, 64)

		likeData := dao.Like{UserId: userId, VideoId: videoId}
		if err := dao.InsertLike(&likeData); err != nil {
			log.Println("Failed to insert likes")
		}
	}
}

func consumerLikeDel(messages <-chan amqp.Delivery) {
	for message := range messages {
		//解析数据
		data := strings.Split(string(message.Body), ":")
		userId, _ := strconv.ParseInt(data[0], 10, 64)
		videoId, _ := strconv.ParseInt(data[1], 10, 64)

		if err := dao.DeleteLike(userId, videoId); err != nil {
			log.Println("Failed to delete likes")
		}
	}
}
