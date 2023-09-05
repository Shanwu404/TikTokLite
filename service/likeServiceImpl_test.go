package service

import (
	"fmt"
	"testing"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/middleware/rabbitmq"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
)

func LikeServiceImplInit() {
	dao.Init()
	redis.InitRedis()
	rabbitmq.Init()
	rabbitmq.InitLikeRabbitMQ()
}

func TestLikeServiceImpl_Like(t *testing.T) {
	LikeServiceImplInit()
	lsi := NewLikeService()
	err := lsi.Like(4321, 1)
	fmt.Println(err)
}

func TestLikeServiceImpl_Unlike(t *testing.T) {
	LikeServiceImplInit()
	lsi := NewLikeService()
	err := lsi.Unlike(4321, 1)
	fmt.Println(err)
}

func TestLikeServiceImpl_GetLikeLists(t *testing.T) {
	LikeServiceImplInit()
	lsi := NewLikeService()
	likes := lsi.GetLikeLists(4321)
	fmt.Println(likes)

}

func TestLikeServiceImpl_IsLike(t *testing.T) {
	LikeServiceImplInit()
	lsi := NewLikeService()
	flag := lsi.IsLike(1, 4321)
	fmt.Println(flag)
}

func TestLikeServiceImpl_CountLikes(t *testing.T) {
	LikeServiceImplInit()
	lsi := NewLikeService()
	cnt := lsi.CountLikes(1)
	fmt.Println(cnt)
}

func TestTotalFavorited(t *testing.T) {
	LikeServiceImplInit()
	lsi := NewLikeService()
	cnt := lsi.TotalFavorited(14)
	fmt.Println(cnt)
}
