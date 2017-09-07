package slss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

const (
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

	fmt.Println("[slss] Uploading lambda function...")
	if err := UploadFunc(apexExecutor); err != nil {
		PrintErrorAndExit(err)
	}

	fmt.Println("[slss] Creating ngrox proxy...")
	proxyAddr, err := StartNgrokProxy(config.Ngrok, ProxyProtoTCP, config.Shadowsocks.LocalPort)
	if err != nil {
		PrintErrorAndExit(err)
	}
	fmt.Println("[slss] Ngrox address: ", proxyAddr)

	fmt.Println("[slss] Starting ss client...")
	localCliCmd, err := StartLocalClient(config, proxyAddr)
	if err != nil {
		PrintErrorAndExit(err)
	}

	defer localCliCmd.Process.Kill()

	requestLambda(apexExecutor, proxyAddr)

	for range time.Tick(interval) {
		requestLambda(apexExecutor, proxyAddr)
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
	_, err := executor.Exec("apex", nil, "deploy", "slss")

	return errors.WithStack(err)
}

func requestLambda(executor *APEXCommandExecutor, proxyAddr string) {
	go func() {
		fmt.Println("[slss] Requesting lambda function")
		if err := RequestRemoteFunc(executor, proxyAddr); err != nil {
			PrintErrorAndExit(err)
		}
	}()
}

// RequestRemoteFunc sends a request to the slss function in AWS lambda
func RequestRemoteFunc(executor *APEXCommandExecutor, proxyAddr string) error {
	proxyHost, proxyPort, err := net.SplitHostPort(proxyAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	lambdaMessage, err := json.Marshal(LambdaShadowSocksConfig{
		Port:      executor.Config.Shadowsocks.ServerPort,
		Method:    executor.Config.Shadowsocks.Method,
		Password:  executor.Config.Shadowsocks.Password,
		ProxyHost: proxyHost,
		ProxyPort: proxyPort,
	})

	if err != nil {
		return errors.WithStack(err)
	}

	_, err = executor.Exec("apex", bytes.NewBufferString(string(lambdaMessage)), "invoke", "slss")

	return errors.WithStack(err)
}
