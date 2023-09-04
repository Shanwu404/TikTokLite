package rabbitmq

import (
	"strconv"
	"strings"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
	"github.com/streadway/amqp"
)

type RelationMQ struct {
	RabbitMQ
	Channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
}

var (
	RabbitMQRelationAdd *RelationMQ
	RabbitMQRelationDel *RelationMQ
)

// NewRelationMQ 创建一个新的操作用户关系的RabbitMQ实例
func NewRelationMQ(queueName string) *RelationMQ {
	relationMQ := &RelationMQ{
		RabbitMQ:  *MyRabbitMQ,
		QueueName: queueName,
	}

	var err error
	relationMQ.Channel, err = relationMQ.Conn.Channel()
	MyRabbitMQ.failOnErr(err, "Failed to get a channel")

	return relationMQ
}

// InitRelationMQ 初始化关系队列
func InitRelationMQ() {
	RabbitMQRelationAdd = NewRelationMQ("Relation Add")
	go RabbitMQRelationAdd.Consumer()
	logger.Infoln("RabbitMQRelationAdd init successfully!")

	RabbitMQRelationDel = NewRelationMQ("Relation Del")
	go RabbitMQRelationDel.Consumer()
	logger.Infoln("RabbitMQRelationDel init successfully!")
}

// Producer 关系队列生产者
func (relationmq *RelationMQ) Producer(msg string) {
	// 1. 申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	_, err := relationmq.Channel.QueueDeclare(
		relationmq.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞
		nil,   // 额外属性
	)
	relationmq.failOnErr(err, "Failed to declare a queue")

	// 2. 发送消息到队列中
	relationmq.Channel.Publish(
		relationmq.Exchange,
		relationmq.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
}

// Consumer 关系队列消费者
func (relationmq *RelationMQ) Consumer() {
	// 1. 申请队列，如果队列不存在会自动创建，如果存在则跳过创建
	_, err := relationmq.Channel.QueueDeclare(
		relationmq.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞
		nil,   // 额外属性
	)
	relationmq.failOnErr(err, "Failed to declare a queue")

	// 2. 接收消息
	msgs, err := relationmq.Channel.Consume(
		relationmq.QueueName,
		"",    // 用来区分多个消费者
		true,  // 是否自动应答
		false, // 是否具有排他性
		false, // 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, // 是否阻塞
		nil,   // 额外属性
	)
	relationmq.failOnErr(err, "Failed to register a consumer")

	forever := make(chan bool)
	switch relationmq.QueueName {
	case "Relation Add":
		go relationmq.consumerFollowAdd(msgs)
	case "Relation Del":
		go relationmq.consumerFollowDel(msgs)
	}
	<-forever

}

// consumerFollowAdd 消费者处理关注关系添加
func (relationmq *RelationMQ) consumerFollowAdd(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		logger.Infof("Received a message: %s", msg.Body)

		// 解析消息体
		params := strings.Split(string(msg.Body), ",")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		followId, _ := strconv.ParseInt(params[1], 10, 64)
		logger.Infoln("this is consumerFollowAdd", userId, followId)

		// 插入新的关注关系
		err := dao.InsertFollow(userId, followId)
		if err == nil {
			logger.Infof("user %d followed user %d successfully", userId, followId)

			// 将新关注关系添加到Redis缓存
			redisFollowKey := utils.RelationFollowKey + strconv.FormatInt(userId, 10)
			redis.RDb.SAdd(redis.Ctx, redisFollowKey, followId)
			// 更新过期时间
			redis.RDb.Expire(redis.Ctx, redisFollowKey, utils.RelationFollowKeyTTL)
		} else {
			logger.Errorln("Failed to insert follow relation: ", err)
		}

	}

}

// consumerFollowDel 消费者处理关注关系删除
func (relationmq *RelationMQ) consumerFollowDel(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		logger.Infof("Received a message: %s", msg.Body)

		// 解析消息体
		params := strings.Split(string(msg.Body), ",")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		followId, _ := strconv.ParseInt(params[1], 10, 64)
		logger.Infoln("this is consumerFollowAdd", userId, followId)

		// 删除关注关系
		err := dao.DeleteFollow(userId, followId)
		if err == nil {
			logger.Infof("user %d unfollowed user %d successfully", userId, followId)

			// 从Redis中移除关注关系
			redisFollowKey := utils.RelationFollowKey + strconv.FormatInt(userId, 10)
			redis.RDb.SRem(redis.Ctx, redisFollowKey, followId)
		} else {
			logger.Errorln("Failed to delete follow relation: ", err)
		}
	}
}
