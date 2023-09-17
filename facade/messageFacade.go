package facade

import (
	"sync"
	"time"

	"github.com/Shanwu404/TikTokLite/service"
)

type MessageInfo struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type MessageActionRequest struct {
	UserId   int64
	ToUserId int64
	Content  string
}

type MessageActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type MessageListRequest struct {
	UserId     int64
	ToUserId   int64
	PreMsgTime int64
}

type MessageListResponse struct {
	Response
	MessageList []MessageInfo `json:"message_list,omitempty"`
}

type MessageFacade struct {
	messageService service.MessageService
}

func NewMessageFacade() *MessageFacade {
	return &MessageFacade{
		messageService: service.NewMessageService(),
	}
}

func (mf *MessageFacade) MessageAction(req MessageActionRequest) MessageActionResponse {
	messageNew := service.MessageParams{
		ToUserId:   req.ToUserId,
		FromUserId: req.UserId,
		Content:    req.Content,
		CreateTime: time.Now(),
	}
	_, code, message := mf.messageService.PublishMessage(messageNew)
	return MessageActionResponse{StatusCode: code, StatusMsg: message}
}

func (mf *MessageFacade) MessageList(req MessageListRequest) MessageListResponse {
	messages := mf.messageService.QueryMessagesByIdsAfter(req.UserId, req.ToUserId, req.PreMsgTime)
	var wg sync.WaitGroup
	wg.Add(len(messages))
	messageList := make([]MessageInfo, len(messages))
	for idx, message := range messages {
		go func(idx int, message service.MessageParams) {
			defer wg.Done()
			messageList[idx] = MessageInfo{
				Id:         message.Id,
				ToUserId:   message.ToUserId,
				FromUserId: message.FromUserId,
				Content:    message.Content,
				CreateTime: message.CreateTime.UnixMilli(),
			}
		}(idx, message)
	}
	wg.Wait()
	return MessageListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		MessageList: messageList,
	}
}
