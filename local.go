package slss

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	commandTemplate = "'%v' | apex invoke slss"
)

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
