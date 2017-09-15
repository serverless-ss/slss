package slss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

	log.Info("[slss] Uploading lambda function...")
	if err := UploadFunc(apexExecutor); err != nil {
		PrintErrorAndExit(err)
	}

	remoteProxyAddrChan := GetProxyAddrChan(config)

	log.Info("[slss] Creating ngrox proxy...")
	proxyAddr, err := StartNgrokProxy(config.Ngrok, ProxyProtoHTTP, config.LocalServerPort)
	if err != nil {
		PrintErrorAndExit(err)
	}
	log.Info("[slss] Ngrox address: ", proxyAddr)

	go func() {
		log.Info("[slss] Request lambda function...")
		if err := RequestRemoteFunc(apexExecutor, proxyAddr); err != nil {
			log.Errorln(err)
		}

		for range time.Tick(interval) {
			log.Info("[slss] Request lambda function...")
			if err := RequestRemoteFunc(apexExecutor, proxyAddr); err != nil {
				log.Errorln(err)
			}
		}
	}()

	var localCliCmd *exec.Cmd
	for remoteProxyAddr := range remoteProxyAddrChan {
		log.Info("[slss] Remote proxy address: ", remoteProxyAddr)

		log.Info("[slss] Restarting local ss client...")
		if localCliCmd != nil {
			if err := localCliCmd.Process.Kill(); err != nil {
				PrintErrorAndExit(err)
			}
		}

		localCliCmd, err = StartLocalClient(config, remoteProxyAddr)
		if err != nil {
			PrintErrorAndExit(err)
		}
		log.Info("[slss] Local ss restarted => localhost:", config.Shadowsocks.LocalPort)
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

// RequestRemoteFunc sends a request to the slss function in AWS lambda
func RequestRemoteFunc(executor *APEXCommandExecutor, proxyAddr string) error {
	lambdaMessage, err := json.Marshal(LambdaShadowSocksConfig{
		Port:       executor.Config.Shadowsocks.ServerPort,
		Method:     executor.Config.Shadowsocks.Method,
		Password:   executor.Config.Shadowsocks.Password,
		ProxyURL:   proxyAddr,
		NgrokToken: executor.Config.Ngrok.AuthToken,
	})

	if err != nil {
		return errors.WithStack(err)
	}

	_, err = executor.Exec("apex", bytes.NewBufferString(string(lambdaMessage)), "invoke", "slss")

	return errors.WithStack(err)
}

// GetProxyAddrChan starts a local server for remote ss server proxy address
func GetProxyAddrChan(config *Config) chan string {
	ch := make(chan string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ssServerAddr := r.URL.Query().Get("ss_server_addr")
		if ssServerAddr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ch <- ssServerAddr
	})

	go func() {
		if err := http.ListenAndServe(":"+config.LocalServerPort, nil); err != nil {
			log.Errorln(err)
		}
	}()
	log.Info("[slss] Local addr chan listen at " + config.LocalServerPort)

	return ch
}
