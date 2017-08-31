package slss

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

const (
	requestCommandTemplate     = "echo '%v' | apex invoke slss"
	deployCommand              = "apex deploy slss"
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

	apexExecutor := &APEXCommandExecutor{Config: config}

	if err := UploadFunc(apexExecutor); err != nil {
		PrintErrorAndExit(err)
	}

	localCliCmd, err := StartLocalClient(config)
	if err != nil {
		PrintErrorAndExit(err)
	}

	defer localCliCmd.Process.Kill()

	for range time.Tick(interval) {
		go func() {
			if err := RequestRemoteFunc(apexExecutor); err != nil {
				PrintErrorAndExit(err)
			}
		}()
	}
}

// StartLocalClient starts a slss client
func StartLocalClient(config *Config) (*exec.Cmd, error) {
	cmd := exec.Command(
		"./bin/shadowsocks_local",
		fmt.Sprintf(localCliServerAddrTemplate, config.Shadowsocks.ServerAddr),
		fmt.Sprintf(localCliServerPortTemplate, config.Shadowsocks.ServerPort),
		fmt.Sprintf(localCliLocalPortTemplate, config.Shadowsocks.LocalPort),
		fmt.Sprintf(localCliPasswordTemplate, config.Shadowsocks.Password),
	)

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "start local client failed")
	}

	return cmd, nil
}

// UploadFunc uploads the slss function to AWS lambda
func UploadFunc(executor *APEXCommandExecutor) error {
	_, err := executor.Exec(deployCommand)

	return errors.Wrap(err, "upload remote function failed")
}

// RequestRemoteFunc sends a request to the slss function in AWS lambda
func RequestRemoteFunc(executor *APEXCommandExecutor) error {
	lambdaMessage, err := json.Marshal(LambdaShadowSocksConfig{
		Addr:     executor.Config.Shadowsocks.ServerAddr,
		Method:   executor.Config.Shadowsocks.Method,
		Password: executor.Config.Shadowsocks.Password,
	})

	if err != nil {
		return errors.Wrap(err, "marshal remote function event failed")
	}

	_, err = executor.Exec(fmt.Sprintf(requestCommandTemplate, lambdaMessage))

	return errors.Wrap(err, "request remote function failed")
}
