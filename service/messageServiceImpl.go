package service

import (
	"log"

	"github.com/Shanwu404/TikTokLite/dao"
)

type MessageServiceImpl struct{}

func NewMessageService() MessageService {
	return &MessageServiceImpl{}
}

func (MessageServiceImpl) QueryMessagesByIds(fromUserId int64, toUserId int64) []MessageParams {
	messages, err := dao.QueryMessagesByIds(fromUserId, toUserId)
	if err != nil {
		log.Println("error:", err.Error())
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
