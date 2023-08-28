package config

import (
	"errors"

	"github.com/BurntSushi/toml"
)

var AppConfig struct {
	HTTPServer HTTPServerConfig
	Database   DatabaseConfig
	OSS        OSSConfig
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

var (
	HTTPServer = &AppConfig.HTTPServer
	Database   = &AppConfig.Database
	OSS        = &AppConfig.OSS
)

func init() {
	_, err := toml.DecodeFile(`config/config.toml`, &AppConfig)
	if err != nil {
		err = errors.Join(errors.New("read config file failed"), err)
		panic(err)
	}
}
