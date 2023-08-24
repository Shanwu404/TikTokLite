package config

import (
	"errors"

	"github.com/BurntSushi/toml"
)

var appConfig struct {
	HTTPServer _HTTPServerConfig
	Database   databaseConfig
	OSS        _OSSConfig
}

type _HTTPServerConfig struct {
	IP   string
	Port int
}

type databaseConfig struct {
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

type _OSSConfig struct {
	CredentialType     string
	CredentialRoleName string
	Endpoint           map[string]string
	BucketName         string
}

func HTTPServer() _HTTPServerConfig {
	return appConfig.HTTPServer
}

func Database() databaseConfig {
	return appConfig.Database
}

func OSS() _OSSConfig {
	return appConfig.OSS
}

func init() {
	_, err := toml.DecodeFile(`config/config.toml`, &appConfig)
	if err != nil {
		err = errors.Join(errors.New("read config file failed"), err)
		panic(err)
	}
}
