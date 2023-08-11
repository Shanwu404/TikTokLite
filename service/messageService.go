package service

import "github.com/Shanwu404/TikTokLite/dao"

type messageService interface {
	// QueryMessagesByIds 根据查询消息列表
	QueryMessagesByIds(fromUserId int64, toUserId int64) []dao.Message

	// PublishMessage 发布消息
	PublishMessage(message dao.Message) (int64, int32, string)
}
