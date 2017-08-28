package slss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
)

const commandTemplate = "'%v' | apex invoke slss"

func requestRemote(config *Config) {
	lambdaMessage, err := json.Marshal(LambdaShadowSocksConfig{
		Addr:     config.Shadowsocks.ServerAddr,
		Method:   config.Shadowsocks.Method,
		Password: config.Shadowsocks.Password,
	})

	if err != nil {
		PrintErrorAndExit(err)
	}

	cmd := exec.Command(fmt.Sprintf(commandTemplate, lambdaMessage))

	var responseMessage bytes.Buffer
	cmd.Stdout = &responseMessage
	cmd.Path = "./lambda"

	if err := cmd.Run(); err != nil {
		PrintErrorAndExit(err)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
}

// StartLocalClient starts a slss client
func StartLocalClient(config *Config) {
	listener, err := net.Listen("tcp", config.Shadowsocks.LocalAddr)
	if err != nil {
		PrintErrorAndExit(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}
