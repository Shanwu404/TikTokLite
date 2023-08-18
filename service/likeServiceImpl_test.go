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
	err := lsi.Like(1000, 1000)
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
	lsi := LikeServiceImpl{}
	likes, err := lsi.GetLikeLists(1000)
	fmt.Println(likes)
	fmt.Println(err)
}

func TestLikeServiceImpl_IsLike(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	flag, err := lsi.IsLike(1000, 1000)
	fmt.Println(flag, err)
}

func TestLikeServiceImpl_LikeCount(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	cnt, err := lsi.LikeCount(1000)
	fmt.Println(cnt, err)
}

func TestTotalFavorited(t *testing.T) {
	LikeServiceImplInit()
	lsi := LikeServiceImpl{}
	cnt := lsi.TotalFavorited(14)
	fmt.Println(cnt)
}
