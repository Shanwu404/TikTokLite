package dao

import (
	"fmt"
	"testing"
)

func TestInsertLike(t *testing.T) {
	Init()
	for i := 0; i < 1; i++ {
		err := InsertLike(&Like{
			UserId:  int64(1002),
			VideoId: int64(1205 + i),
		})
		fmt.Printf("%v", err)
	}

}

func TestDeleteLike(t *testing.T) {
	Init()
	for i := 4; i < 5; i++ {
		err := DeleteLike(int64(1001), int64(1200+i))
		fmt.Printf("%v", err)
	}
}

func TestGetLikeVideoIdList(t *testing.T) {
	Init()
	list, err := GetLikeVideoIdList(1002)
	fmt.Println(list)
	fmt.Println(err)
}

func TestCountLikes(t *testing.T) {
	Init()
	list, err := CountLikes(1200)
	fmt.Println(list)
	fmt.Println(err)
}
