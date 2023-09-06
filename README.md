# TikTokLite

![GitHub Repo stars](https://img.shields.io/github/stars/Shanwu404/TikTokLite)
![GitHub watchers](https://img.shields.io/github/watchers/Shanwu404/TikTokLite)
![GitHub forks](https://img.shields.io/github/forks/Shanwu404/TikTokLite)
![GitHub contributors](https://img.shields.io/github/contributors/Shanwu404/TikTokLite)

基于 Gin 和 GORM 的仿抖音后端，采用 MVC 架构搭建，引入 Redis 和 RabbitMQ 进行优化，实现了低负载和高性能。
<div align="center">
<img src="images/logo.png" alt="logo" width="472" height="328"/>
</div>


## 目录

- [使用说明](#使用说明)
- [项目依赖](#项目依赖)
- [项目配置](#项目配置)
  - [配置文件](#配置文件)
  - [MySQL配置](#MySQL配置)
  - [Redis配置](#Redis配置)
  - [RabbitMQ配置](#RabbitMQ配置)


## 使用说明

项目整体架构：

```
TikTokLite 
├── /config/ 配置文件包
├── /controller/ 控制器包
├── /dao/ 数据库访问
├── /images/ 图片引用
├── /middleware/ 中间件
│   ├── ffmpeg/ 视频截图
│   ├── jwt/ 鉴权
│   ├── rabbitmq/ 消息队列
│   ├── redis/ 缓存
├── /service/ 服务层
├── /utils/ 工具
├── .gitignore
├── /go.mod/
├── LICENSE
├── main.go
├── README.md
└── router.go
```

项目运行：

```go
go run main.go router.go
```

## 项目依赖

TikTokLite 项目依赖如下：

```
module github.com/Shanwu404/TikTokLite

go 1.20

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/aliyun/aliyun-oss-go-sdk v2.2.8+incompatible
	github.com/aliyun/credentials-go v1.3.1
	github.com/brianvoe/gofakeit/v6 v6.23.1
	github.com/gin-contrib/pprof v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/importcjj/sensitive v0.0.0-20200106142752-42d1c505be7b
	github.com/streadway/amqp v1.1.0
	go.uber.org/zap v1.25.0
	golang.org/x/crypto v0.9.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.2
)
```

依赖安装：

```bash
go mod download
```

## 项目配置

### 配置文件

本项目的所有配置都在 .toml 文件中，出于安全考虑并未上传到项目中，该文件包含内容如下：

```
[HTTPServer]
IP = "" // 服务器 IP
Port =  // 服务器端口号
[Database]
IP = "" // 数据库 IP
Port =  // 数据库端口号
Account = "" // 数据库用户名
Password = "" // 数据库密码
DatabaseName = "" // 数据库名
Protocol = ""
Charset = ""
ParseTime = 
TimeZone = ""
[OSS]
CredentialType = ""
CredentialRoleName = ""
Endpoint = { Internal = "oss-cn-beijing-internal.aliyuncs.com", External = "oss-cn-beijing.aliyuncs.com" }
BucketName = ""
[Redis]
RedisHost = "" // Redis IP
RedisPort =  // Redis 端口号
RedisPassword = "" // Redis 密码
[Rabbitmq]
RabbitmqHost = "" // 消息队列 IP
RabbitmqPort =  // 消息队列端口号
RabbitmqUsername = "" // 消息队列用户名
RabbitmqPassword = "" // 消息队列密码
```

### MySQL配置

MySQL 的安装配置比较简单，此处略过⏭

### Redis配置

**安装**

官网下载安装包解压即可

**启动/测试**

```bash
redis-server
ping # 返回 PONG
```

### RabbitMQ配置

**安装**

```bash
# Linux
apt-get install rabbitmq-server
# Mac
brew intsall rabbitmq-server
```

**启动**

```bash
rabbitmq-server
```
