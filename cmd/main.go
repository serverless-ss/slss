package main

import (
	"flag"
	"fmt"
	"os"

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
		printErrorAndExit(err)
	}

	fmt.Printf("config: %+v\n", config)
}

func printErrorAndExit(err error) {
	fmt.Printf("%+v\n", err)
	os.Exit(-1)
}
