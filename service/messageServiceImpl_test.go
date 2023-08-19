package service

import (
	"testing"

	"github.com/Shanwu404/TikTokLite/dao"
)

func MessageServiceImplInit() {
	dao.Init()
}

func TestMessageServiceImpl_QueryMessagesByIds(t *testing.T) {
	// MessageServiceImplInit()
	// msi := MessageServiceImpl{}
	// messages := msi.QueryMessagesByIds(0, 1)
	// fmt.Println(messages)
}

func TestMessageServiceImpl_PublishMessage(t *testing.T) {
	// MessageServiceImplInit()
	// msi := MessageServiceImpl{}
	// message := dao.Message{
	// 	FromUserId: 0,
	// 	ToUserId:   1,
	// 	Content:    "test",
	// 	CreateTime: time.Now(),
	// }
	// id, code, newMessage := msi.PublishMessage(message)
	// fmt.Println(id, code, newMessage)
}
