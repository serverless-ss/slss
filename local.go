package slss

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

const (
	requestCommandTemplate     = "echo '%v' | apex invoke slss"
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

	proxyAddr, err := StartNgrokProxy(config.Ngrok, ProxyProtoTCP, config.Shadowsocks.LocalPort)
	if err != nil {
		PrintErrorAndExit(err)
	}

	localCliCmd, err := StartLocalClient(config, proxyAddr)
	if err != nil {
		PrintErrorAndExit(err)
	}

	defer localCliCmd.Process.Kill()

	for range time.Tick(interval) {
		go func() {
			if err := RequestRemoteFunc(apexExecutor, proxyAddr); err != nil {
				PrintErrorAndExit(err)
			}
		}()
	}
}

// StartLocalClient starts a slss client
func StartLocalClient(config *Config, proxyAddr string) (*exec.Cmd, error) {
	host, port, err := net.SplitHostPort(proxyAddr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cmd := exec.Command(
		"./bin/shadowsocks_local",
		fmt.Sprintf(localCliServerAddrTemplate, host),
		fmt.Sprintf(localCliServerPortTemplate, port),
		fmt.Sprintf(localCliLocalPortTemplate, config.Shadowsocks.LocalPort),
		fmt.Sprintf(localCliPasswordTemplate, config.Shadowsocks.Password),
	)

	if err := cmd.Start(); err != nil {
		return nil, errors.WithStack(err)
	}

	return cmd, nil
}

// UploadFunc uploads the slss function to AWS lambda
func UploadFunc(executor *APEXCommandExecutor) error {
	_, err := executor.Exec("apex", "deploy", "slss")

	return errors.WithStack(err)
}

// RequestRemoteFunc sends a request to the slss function in AWS lambda
func RequestRemoteFunc(executor *APEXCommandExecutor, proxyAddr string) error {
	proxyHost, proxyPort, err := net.SplitHostPort(proxyAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	lambdaMessage, err := json.Marshal(LambdaShadowSocksConfig{
		Addr:      executor.Config.Shadowsocks.ServerAddr,
		Method:    executor.Config.Shadowsocks.Method,
		Password:  executor.Config.Shadowsocks.Password,
		ProxyHost: proxyHost,
		ProxyPort: proxyPort,
	})

	if err != nil {
		return errors.WithStack(err)
	}

	_, err = executor.Exec(fmt.Sprintf(requestCommandTemplate, lambdaMessage))

	return errors.WithStack(err)
}
