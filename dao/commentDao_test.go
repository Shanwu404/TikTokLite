package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestQueryCommentsByVideoId(t *testing.T) {
	Init()
	comments, err := QueryCommentsByVideoId(0)
	fmt.Println(comments)
	fmt.Println(err)
}

func TestInsertComment(t *testing.T) {
	Init()
	comment := Comment{
		UserId:     1000,
		VideoId:    1000,
		Content:    "hello",
		CreateDate: time.Now(),
	}
	newComment, err := InsertComment(comment)
	fmt.Println(newComment)
	fmt.Println(err)
}

func TestDeleteComment(t *testing.T) {
	Init()
	flag := DeleteComment(3)
	fmt.Println(flag)
}
