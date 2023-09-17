package controller

import (
	"net/http"

	"github.com/Shanwu404/TikTokLite/facade"
	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageFacade facade.MessageFacade
}

func NewMessageController() *MessageController {
	return &MessageController{
		messageFacade: *facade.NewMessageFacade(),
	}
}

// MessageAction POST /douyin/message/action/ 发送消息
func (ms *MessageController) MessageAction(c *gin.Context) {
	req, valid := MessageActionParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, facade.MessageActionResponse{StatusCode: -1, StatusMsg: "Invalid Request."})
		return
	}
	c.JSON(http.StatusOK, ms.messageFacade.MessageAction(req))
	return
}

// MessageList GET /douyin/message/chat/ 聊天记录
func (ms *MessageController) MessageList(c *gin.Context) {
	req, valid := MessageListParseAndValidateParams(c)
	if !valid {
		c.JSON(http.StatusBadRequest, facade.MessageListResponse{
			Response: facade.Response{StatusCode: -1, StatusMsg: "Invalid Request."},
		})
		return
	}
	c.JSON(http.StatusOK, ms.messageFacade.MessageList(req))
	return
}
