package service

import (
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
	if len(messages) > 0 {
		messages = messages[1:]
	}
	results := make([]MessageParams, 0, len(messages))
	for i := range messages {
		results = append(results, MessageParams(messages[i]))
	}
	log.Println("Query messages successfully!")
	return results
}

func (MessageServiceImpl) PublishMessage(message MessageParams) (int64, int32, string) {
	messageNew, err := dao.InsertMessage(dao.Message(message))
	if err != nil {
		return -1, 1, "Publish message failed!"
	}
	return messageNew.Id, 0, "Publish message successfully!"
}
