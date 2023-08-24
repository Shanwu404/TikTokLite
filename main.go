package main

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
)

func main() {
	r := NewRouter()
	initDeps()
	//utils.FakeComments(10)

	err := r.Run() // listen and serve on listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}

func initDeps() {
	//初始化数据库连接
	dao.Init()
	redis.InitRedis()
}
