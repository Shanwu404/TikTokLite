package service

import (
	"fmt"
	"testing"

	"github.com/Shanwu404/TikTokLite/dao"
)

func LikeServiceImplInit() {
	dao.Init()
}

func TestLikeServiceImpl_Like(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	err := lsi.Like(4321, 1)
	fmt.Println(err)
}

func TestLikeServiceImpl_Unlike(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	err := lsi.Unlike(1000, 1000)
	fmt.Println(err)
}

func TestLikeServiceImpl_GetLikeLists(t *testing.T) {
	LikeServiceImplInit()
	lsi := NewLikeService()
	// lsi := LikeServiceImpl{}
	likes := lsi.GetLikeLists(4321)
	fmt.Println(likes)

}

func TestLikeServiceImpl_IsLike(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	flag := lsi.IsLike(1000, 1000)
	fmt.Println(flag)
}

func TestLikeServiceImpl_CountLikes(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	cnt := lsi.CountLikes(1000)
	fmt.Println(cnt)
}

func TestTotalFavorited(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	cnt := lsi.TotalFavorited(14)
	fmt.Println(cnt)
}
