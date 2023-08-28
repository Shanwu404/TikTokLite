package main

import (
	"fmt"

	"github.com/Shanwu404/TikTokLite/config"
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
)

func main() {
	r := NewRouter()
	initDeps()
	//utils.FakeComments(10)

	port := config.HTTPServer.Port
	err := r.Run(fmt.Sprintf(":%d", port)) // listen and serve on listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}

func initDeps() {
	//初始化数据库连接
	dao.Init()
	redis.InitRedis()
	utils.InitWordsFilter()
}
