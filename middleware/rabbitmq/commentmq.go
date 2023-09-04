package rabbitmq

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/streadway/amqp"
)

type CommentMQ struct {
	RabbitMQ
	Channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
}

var CommentDel *CommentMQ

// NewCommentMQ 获取Comment的消息队列
func NewCommentMQ(queueName string) *CommentMQ {
	commentMQ := &CommentMQ{
		RabbitMQ:  *MyRabbitMQ,
		QueueName: queueName,
	}
	var err error
	commentMQ.Channel, err = commentMQ.Conn.Channel()
	MyRabbitMQ.failOnErr(err, "Failed to get channel!")
	return commentMQ
}

func InitCommentMQ() {
	CommentDel = NewCommentMQ("Comment Del")
	go CommentDel.Consumer()
	log.Println("RabbitMQCommentDel init successfully!")
}

// Producer 生产
func (c *CommentMQ) Producer(message string) {
	_, err := c.Channel.QueueDeclare(
		c.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
	}
	err = c.Channel.Publish(
		c.Exchange,
		c.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Println(err.Error())
	}
}

// Consumer 消费
func (c *CommentMQ) Consumer() {
	_, err := c.Channel.QueueDeclare(
		c.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
	}
	messages, err := c.Channel.Consume(
		c.QueueName,
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
	go consumerCommentDel(messages)
	<-forever
}

func consumerCommentDel(messages <-chan amqp.Delivery) {
	for message := range messages {
		id := fmt.Sprintf("%s", message.Body)
		commentId, _ := strconv.Atoi(id)
		flag := dao.DeleteComment(int64(commentId))
		if !flag {
			log.Println("Comment delete failed!")
		}
	}
}
