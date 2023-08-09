package main

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	initDeps()
	//utils.FakeComments(10)

	initRouter(r)

	err := r.Run() // listen and serve on listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}

func initDeps() {
	//初始化数据库连接
	dao.Init()
}
