package config

import (
	"errors"

	"github.com/BurntSushi/toml"
)

var Redis_host = "127.0.0.1"
var Redis_port = "6379"
var Redis_password = ""

var AppConfig struct {
	HTTPServer HTTPServerConfig
	Database   DatabaseConfig
	OSS        OSSConfig
	Redis      RedisConfig
}

type HTTPServerConfig struct {
	IP   string
	Port int
}

type DatabaseConfig struct {
	IP           string
	Port         int
	Account      string
	Password     string
	Protocol     string
	DatabaseName string
	Charset      string
	ParseTime    bool
	TimeZone     string
}

type OSSConfig struct {
	CredentialType     string
	CredentialRoleName string
	Endpoint           map[string]string
	BucketName         string
}

type RedisConfig struct {
	Redis_host     string
	Redis_port     int
	Redis_password string
}

var (
	HTTPServer = &AppConfig.HTTPServer
	Database   = &AppConfig.Database
	OSS        = &AppConfig.OSS
	Redis      = &AppConfig.Redis
)

func init() {
	_, err := toml.DecodeFile(`../config/config.toml`, &AppConfig)
	if err != nil {
		err = errors.Join(errors.New("read config file failed"), err)
		panic(err)
	}
}
