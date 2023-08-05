package utils

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/brianvoe/gofakeit/v6"
	"time"
)

func FakeComments(num int) {
	gofakeit.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		comment := dao.Comment{}
		comment.UserId = int64(i)
		comment.VideoId = int64(i)
		comment.Content = gofakeit.Sentence(20)
		comment.CreateDate = gofakeit.Date()
		_, err := dao.InsertComment(comment)
		if err != nil {
			return
		}
	}
}
