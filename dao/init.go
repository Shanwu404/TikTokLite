package dao

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shanwu404/TikTokLite/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // 彩色打印
		},
	)

	var (
		myConfig     = config.Database()
		account      = myConfig.Account
		password     = myConfig.Password
		ip           = myConfig.IP
		port         = myConfig.Port
		databaseName = myConfig.DatabaseName
		protocol     = myConfig.Protocol
		charset      = myConfig.Charset
		parsetime    = myConfig.ParseTime
		timeZone     = myConfig.TimeZone
	)

	dsn := fmt.Sprintf(
		"%s:%s@%s(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		account, password, protocol, ip, port, databaseName, charset, parsetime, timeZone)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panicln("err:", err.Error())
	} else {
		log.Println("Database is connected successfully.")
	}
	db.Begin()
}
