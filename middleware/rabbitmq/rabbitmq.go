package rabbitmq

import (
	"fmt"
	"strconv"

	"github.com/Shanwu404/TikTokLite/config"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Url  string
}

var MyRabbitMQ *RabbitMQ

func Init() {
	MyRabbitMQ = &RabbitMQ{
		Url: "amqp://" + config.Rabbitmq().RabbitmqUsername + ":" + config.Rabbitmq().RabbitmqPassword + "@" + config.Rabbitmq().RabbitmqHost + ":" + strconv.Itoa(config.Rabbitmq().RabbitmqPort) + "/",
	}
	var err error
	MyRabbitMQ.Conn, err = amqp.Dial(MyRabbitMQ.Url)
	MyRabbitMQ.failOnErr(err, "Create connection failed!")
	logger.Infoln("Rabbitmq has connected!")
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		logger.Errorln("%s:%s\n", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.Conn.Close()
}
