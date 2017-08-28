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
	flag.StringVar(&configFilePath, "c", "", "path to the configuration file")
	flag.Parse()

	config, err := slss.LoadConfig(configFilePath)
	if err != nil {
		slss.PrintErrorAndExit(err)
	}

	fmt.Printf("config: %+v\n", config)

	go slss.StartLocalClient(config)
}
