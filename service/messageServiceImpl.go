package service

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"log"
)

type MessageServiceImpl struct{}

func (MessageServiceImpl) QueryMessagesByIds(fromUserId int64, toUserId int64) []dao.Message {
	messages, err := dao.QueryMessagesByIds(fromUserId, toUserId)
	if err != nil {
		log.Println("error:", err.Error())
	}
	log.Println("Query messages successfully!")
	return messages
}

func (MessageServiceImpl) PublishMessage(message dao.Message) (int64, int32, string) {
	message, err := dao.InsertMessage(message)
	if err != nil {
		return -1, 1, "Publish message failed!"
	}
	return message.Id, 0, "Publish message successfully!"
}
