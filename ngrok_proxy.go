package slss

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Proxy types
const (
	ProxyProtoHTTP  = "http"
	ProxyProtoHTTPS = "https"
	ProxyProtoTCP   = "tcp"
)

const ngrokBinPath = "./bin/ngrok"

// StartNgrokProxy starts the ngrok proxy
func StartNgrokProxy(config *ngrokConfig, protoType string, port string) (string, error) {
	if err := authNgrok(config.AuthToken); err != nil {
		return "", errors.WithStack(err)
	}

	return start(protoType, port)
}

func authNgrok(authToken string) error {
	cmd := exec.Command(ngrokBinPath, "authtoken", authToken)
	return cmd.Run()
}

func start(proxyType string, port string) (string, error) {
	var responseMessage bytes.Buffer

	cmd := exec.Command(ngrokBinPath, proxyType, port, "-log=stdout", "--log-level=debug", "--region=ap")
	cmd.Stdout = &responseMessage

	if err := cmd.Start(); err != nil {
		return "", errors.WithStack(err)
	}

	go cmd.Wait()

	for range time.Tick(time.Second) {
		output := responseMessage.String()

		if !strings.Contains(output, "tcp://") {
			continue
		}

		i := strings.LastIndex(output, "tcp://")

		return output[i+len("tcp://") : i+strings.Index(output[i:], " ")], nil
	}

	return "", errors.New("unreachable")
}
