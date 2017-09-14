package main

import (
	"flag"

	"github.com/serverless-ss/slss"
	log "github.com/sirupsen/logrus"
)

var (
	configFilePath string
)

func main() {
	flag.StringVar(&configFilePath, "c", "./config.json", "path to the configuration file")
	flag.Parse()

	config, err := slss.LoadConfig(configFilePath)
	if err != nil {
		slss.PrintErrorAndExit(err)
	}

	funcConfig, err := slss.LoadFuncConfig("./lambda/functions/slss/function.json")
	if err != nil {
		slss.PrintErrorAndExit(err)
	}

	log.WithFields(log.Fields{
		"AWS.access_key_id":       config.AWS.AccessKeyID,
		"AWS.secret_access_key":   config.AWS.AccessKeyID,
		"AWS.region":              config.AWS.Region,
		"shadowsocks.server_addr": config.Shadowsocks.ServerAddr,
		"shadowsocks.server_port": config.Shadowsocks.ServerPort,
		"shadowsocks.local_port":  config.Shadowsocks.LocalPort,
		"shadowsocks.timeout":     config.Shadowsocks.Timeout,
		"shadowsocks.method":      config.Shadowsocks.Method,
		"shadowsocks.password":    config.Shadowsocks.Password,
		"ngrok.auth_token":        config.Ngrok.AuthToken,
	}).Info("[slss] Config:")
	log.WithFields(log.Fields{
		"name":        funcConfig.Name,
		"description": funcConfig.Description,
		"runtime":     funcConfig.Runtime,
		"memory":      funcConfig.Memory,
		"timeout":     funcConfig.Timeout,
	}).Info("[slss] Lambda function config:")

	slss.Init(config, funcConfig)
}
