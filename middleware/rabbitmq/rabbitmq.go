package rabbitmq

import (
	"fmt"
	"log"

	"github.com/Shanwu404/TikTokLite/config"
	"github.com/streadway/amqp"
)

var MQURL = "amqp://" + config.RabbitMQ_username + ":" + config.RabbitMQ_passsword + "@" + config.RabbitMQ_IP + ":" + config.RabbitMQ_host + "/"

type RabbitMQ struct {
	Conn *amqp.Connection
	//MQ链接字符串
	Mqurl string
}

var MyRabbitMQ *RabbitMQ

// Init 创建结构体实例
func Init() {
	MyRabbitMQ = &RabbitMQ{
		Mqurl: MQURL,
	}
	var err error
	MyRabbitMQ.Conn, err = amqp.Dial(MyRabbitMQ.Mqurl)
	MyRabbitMQ.failOnErr(err, "Create connection failed!")
	log.Println("rabbitmq has connected!")
}

// failOnErr 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Printf("%s:%s\n", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.Conn.Close()
}
