package redis

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/Shanwu404/TikTokLite/config"
	"github.com/go-redis/redis/v8"
)

var mutex sync.Mutex
var RDb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RDb = redis.NewClient(&redis.Options{
		Addr:     config.Redis().RedisHost + ":" + strconv.Itoa(config.Redis().RedisPort),
		Password: config.Redis().RedisPassword,
		DB:       0,
	})
	_, err := RDb.Ping(Ctx).Result()
	if err != nil {
		log.Println("err:", err.Error())
		return
	}
	log.Println("Redis has connected!")
}

func Lock(key string, value string) bool {
	mutex.Lock() // 保证程序不存在并发冲突问题
	defer mutex.Unlock()
	ret, err := RDb.SetNX(Ctx, key, value, time.Second).Result()
	if err != nil {
		log.Println("Lock error:", err.Error())
		return ret
	}
	return ret
}

func Unlock(key string) bool {
	err := RDb.Del(Ctx, key).Err()
	if err != nil {
		log.Println("Unlock error:", err.Error())
		return false
	}
	return true
}

func RandomTime() time.Duration {
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Int63n(25)) * time.Hour // 设置随机过期时间
}
