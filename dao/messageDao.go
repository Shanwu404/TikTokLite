package dao

import (
	"log"
	"time"
)

type Message struct {
	Id         int64
	ToUserId   int64
	FromUserId int64
	Content    string
	CreateTime time.Time
}

// QueryMessagesByIds 根据fromUserId和toUserId获取所有消息记录
func QueryMessagesByIds(fromUserId int64, toUserId int64) ([]Message, error) {
	var messages []Message
	if err := db.Where("to_user_id = ? AND from_user_id = ?", toUserId, fromUserId).Or("to_user_id = ? AND from_user_id = ?", fromUserId, toUserId).Find(&messages).Error; err != nil {
		log.Println(err)
		return messages, err
	}
	return messages, nil
}

// InsertMessage 插入消息
func InsertMessage(message Message) (Message, error) {
	if err := db.Model(Message{}).Create(&message).Error; err != nil {
		log.Println(err)
		return Message{}, err
	}
	return message, nil
}