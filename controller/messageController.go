package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Shanwu404/TikTokLite/utils/validation"

	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type MessageListResponse struct {
	Response
	MessageList []MessageInfo `json:"message_list,omitempty"`
}

type MessageController struct {
	messageService service.MessageService
}

func NewMessageController() *MessageController {
	return &MessageController{
		messageService: service.NewMessageService(),
	}
}

// MessageAction POST /douyin/message/action/ 发送消息
func (ms *MessageController) MessageAction(c *gin.Context) {
	req, valid := validation.MessageActionParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, douyinPublishActionResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}
	message := service.MessageParams{
		ToUserId:   req.ToUserId,
		FromUserId: req.UserId,
		Content:    req.Content,
		CreateTime: time.Now(),
	}
	_, code, messgae := ms.messageService.PublishMessage(message)
	c.JSON(http.StatusOK, Response{StatusCode: code, StatusMsg: messgae})
	return
}

// MessageList GET /douyin/message/chat/ 聊天记录
func (ms *MessageController) MessageList(c *gin.Context) {
	req, valid := validation.MessageListParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, douyinPublishActionResponse{
			Response: Response{-1, "Invalid Request."},
		})
		return
	}
	// 获取此后的消息
	preMsgTimeStr := c.Query("pre_msg_time")
	preMsgTime, _ := strconv.ParseInt(preMsgTimeStr, 10, 64)
	messages := ms.messageService.QueryMessagesByIdsAfter(req.UserId, req.ToUserId, preMsgTime)
	messageList := make([]MessageInfo, 0, len(messages))
	for _, message := range messages {
		messageList = append(messageList, MessageInfo{
			Id:         message.Id,
			ToUserId:   message.ToUserId,
			FromUserId: message.FromUserId,
			Content:    message.Content,
			CreateTime: message.CreateTime.UnixMilli(),
		})
	}
	c.JSON(http.StatusOK, MessageListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		MessageList: messageList,
	})
}
