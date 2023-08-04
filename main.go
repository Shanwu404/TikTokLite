package main

import (
	"github.com/Shanwu404/TikTokLite/dao"
)

func main() {
	initDeps()
}

func initDeps() {
	//初始化数据库连接
	dao.Init()
}
