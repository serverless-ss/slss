package main

import (
	"flag"
	"fmt"

	"github.com/serverless-ss/slss"
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

	fmt.Printf("config: %+v\n", config)
	fmt.Printf("function config: %+v\n", funcConfig)

	slss.Init(config, funcConfig)
}
