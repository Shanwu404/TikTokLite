package service

import "time"

type MessageParams struct {
	Id         int64     `json:"id"`
	ToUserId   int64     `json:"to_user_id"`
	FromUserId int64     `json:"from_user_id"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

type MessageService interface {
	// QueryMessagesByIds 根据查询消息列表
	QueryMessagesByIdsAfter(fromUserId int64, toUserId int64, timestamp int64) []MessageParams

	// PublishMessage 发布消息
	PublishMessage(message MessageParams) (int64, int32, string)
}
