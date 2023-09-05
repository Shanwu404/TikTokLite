package main

import (
	"fmt"

	"github.com/Shanwu404/TikTokLite/middleware/rabbitmq"

	"github.com/Shanwu404/TikTokLite/config"
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
)

func main() {
	defer logger.Sync()
	r := NewRouter()
	initDeps()

	myConfig := config.HTTPServer()
	err := r.Run(fmt.Sprintf(":%d", myConfig.Port)) // listen and serve on listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return
	}
}

func initDeps() {
	//初始化数据库连接
	dao.Init()
	rabbitmq.Init()
	rabbitmq.InitCommentMQ()
	rabbitmq.InitRelationMQ()
	rabbitmq.InitLikeRabbitMQ()
	redis.InitRedis()
	utils.InitWordsFilter()
}
