package dao

import (
	"github.com/Shanwu404/TikTokLite/log/logger"
	"time"
)

type Message struct {
	Id         int64     `json:"id"`
	ToUserId   int64     `json:"to_user_id"`
	FromUserId int64     `json:"from_user_id"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

// QueryMessagesByIds 根据fromUserId和toUserId获取所有消息记录
func QueryMessagesByIdsAfter(fromUserId int64, toUserId int64, createdTime time.Time) ([]Message, error) {
	var messages []Message
	if err := db.
		Where("to_user_id = ? AND from_user_id = ? AND create_time >= ?", toUserId, fromUserId, createdTime).
		Or("to_user_id = ? AND from_user_id = ? AND create_time >= ?", fromUserId, toUserId, createdTime).
		Order("create_time ASC").
		Find(&messages).Error; err != nil {
		logger.Errorln(err)
		return messages, err
	}
	return messages, nil
}

// InsertMessage 插入消息
func InsertMessage(message Message) (Message, error) {
	if err := db.Model(Message{}).Create(&message).Error; err != nil {
		logger.Errorln(err)
		return Message{}, err
	}
	return message, nil
}
