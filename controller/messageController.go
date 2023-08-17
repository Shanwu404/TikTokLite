package controller

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type MessageListResponse struct {
	Response
	MessageList []dao.Message
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
	actionType := c.Query("action_type")
	msi := service.MessageServiceImpl{}
	if actionType == "1" {
		// 获取当前用户
		userId := c.GetInt64("id")
		// 获取接受用户
		id := c.Query("to_user_id")
		toUserId, _ := strconv.ParseInt(id, 10, 64)
		// 获取内容
		content := c.Query("content")
		message := dao.Message{
			ToUserId:   toUserId,
			FromUserId: userId,
			Content:    content,
			CreateTime: time.Now(),
		}
		_, code, messgae := msi.PublishMessage(message)
		c.JSON(http.StatusOK, Response{StatusCode: code, StatusMsg: messgae})
		return
	}
}

// MessageList GET /douyin/message/chat/ 聊天记录
func (ms *MessageController) MessageList(c *gin.Context) {
	msi := service.MessageServiceImpl{}
	// 获取当前用户
	userId := c.GetInt64("id")
	// 获取接受用户
	id := c.Query("to_user_id")
	toUserId, _ := strconv.ParseInt(id, 10, 64)
	messages := msi.QueryMessagesByIds(userId, toUserId)
	c.JSON(http.StatusOK, MessageListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		MessageList: messages,
	})
	return
}
