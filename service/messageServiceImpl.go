package service

import (
	"github.com/Shanwu404/TikTokLite/log/logger"
	"log"
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
)

type MessageServiceImpl struct{}

func NewMessageService() MessageService {
	return &MessageServiceImpl{}
}

func (MessageServiceImpl) QueryMessagesByIdsAfter(fromUserId int64, toUserId int64, timestamp int64) []MessageParams {
	msgTime := time.UnixMilli(timestamp)
	messages, err := dao.QueryMessagesByIdsAfter(fromUserId, toUserId, msgTime)
	if err != nil {
		log.Println("error:", err.Error())
	}
	//这样的处理逻辑放在实际生产环境中其实不是好的解决方案，比如发送时对方刚好也发消息，并且对方消息实际插入数据库时间比”我“的早
	//主要问题还是这块的前端逻辑不好
	if len(messages) > 0 && msgTime.UnixMilli() != 0 {
		messages = messages[1:]
	}
	results := make([]MessageParams, 0, len(messages))
	for i := range messages {
		results = append(results, MessageParams(messages[i]))
	}
	logger.Infoln("Query messages successfully!")
	return results
}

func (MessageServiceImpl) PublishMessage(message MessageParams) (int64, int32, string) {
	messageNew, err := dao.InsertMessage(dao.Message(message))
	if err != nil {
		return -1, 1, "Publish message failed!"
	}
	return messageNew.Id, 0, "Publish message successfully!"
}
