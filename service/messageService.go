package service

type MessageParams struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type MessageService interface {
	// QueryMessagesByIds 根据查询消息列表
	QueryMessagesByIds(fromUserId int64, toUserId int64) []MessageParams

	// PublishMessage 发布消息
	PublishMessage(message MessageParams) (int64, int32, string)
}
