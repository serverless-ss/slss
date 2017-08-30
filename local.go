package slss

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

const (
	commandTemplate            = "'%v' | apex invoke slss"
	localCliServerAddrTemplate = "-s %v"
	localCliPasswordTemplate   = "-k %v"
	localCliServerPortTemplate = "-p %v"
	localCliLocalPortTemplate  = "-l %v"
)

// Init starts the slss
func Init(config *Config, funcConfig *FuncConfig) {
	interval, err := time.ParseDuration(fmt.Sprintf("%vs", funcConfig.Timeout-10))
	if err != nil {
		PrintErrorAndExit(err)
	}

	localCliCmd := StartLocalClient(config)
	defer localCliCmd.Process.Kill()

	for range time.Tick(interval) {
		go requestRemote(config)
	}
}

func requestRemote(config *Config) {
	lambdaMessage, err := json.Marshal(LambdaShadowSocksConfig{
		Addr:     config.Shadowsocks.ServerAddr,
		Method:   config.Shadowsocks.Method,
		Password: config.Shadowsocks.Password,
	})

	if err != nil {
		PrintErrorAndExit(err)
	}

	executor := &APEXCommandExecutor{Config: config}

	if _, err = executor.Exec(fmt.Sprintf(commandTemplate, lambdaMessage)); err != nil {
		PrintErrorAndExit(err)
	}
}

// StartLocalClient starts a slss client
func StartLocalClient(config *Config) *exec.Cmd {
	cmd := exec.Command(
		"./bin/shadowsocks_local",
		fmt.Sprintf(localCliServerAddrTemplate, config.Shadowsocks.ServerAddr),
		fmt.Sprintf(localCliServerPortTemplate, config.Shadowsocks.ServerPort),
		fmt.Sprintf(localCliLocalPortTemplate, config.Shadowsocks.LocalPort),
		fmt.Sprintf(localCliPasswordTemplate, config.Shadowsocks.Password),
	)

	if err := cmd.Start(); err != nil {
		PrintErrorAndExit(err)
	}

	return cmd
}
