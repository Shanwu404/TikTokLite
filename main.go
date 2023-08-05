package main

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	initDeps()
	//utils.FakeComments(10)
	err := r.Run()
	if err != nil {
		return
	}
}

func initDeps() {
	//初始化数据库连接
	dao.Init()
}
