package config

import (
	"errors"
	"flag"
	"fmt"

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

// type logConfig struct {
// 	LogRootPath string
// 	LogLevel    []string
// }

func HTTPServer() _HTTPServerConfig {
	return appConfig.HTTPServer
}

func Database() databaseConfig {
	return appConfig.Database
}

func OSS() _OSSConfig {
	return appConfig.OSS
}

// func Log() logConfig {
// 	return appConfig.log
// }

// 无论被import多少次init()都只执行一次
func init() {
	configFilePath := "config/config_debug.toml"
	mode := flag.String("mode", "debug", `"debug" or "release"`)
	switch *mode {
	case "release":
		configFilePath = `config/config_` + *mode + `.toml`
	default:
	}
	_, err := toml.DecodeFile(configFilePath, &appConfig)
	fmt.Println("\n\n   111", appConfig)
	if err != nil {
		err = errors.Join(errors.New("read config file failed"), err)
		panic(err)
	}
}
