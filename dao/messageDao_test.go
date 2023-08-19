package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestQueryMessagesByIds(t *testing.T) {
	// Init()
	// messages, err := QueryMessagesByIds(0, 1)
	// fmt.Println(messages)
	// fmt.Println(err)
}

func TestInsertMessage(t *testing.T) {
	Init()
	message := Message{
		FromUserId: 0,
		ToUserId:   1,
		Content:    "test",
		CreateTime: time.Now(),
	}
	newMessage, err := InsertMessage(message)
	fmt.Println(newMessage)
	fmt.Println(err)
}
