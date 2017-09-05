package slss

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Amazon AWS configuration
type awsConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Region          string `json:"region"`
}

// Shadowsocks configuration
type shadowsocksConfig struct {
	ServerAddr string `json:"server_addr"`
	ServerPort string `json:"server_port"`
	LocalPort  string `json:"local_port"`
	Timeout    int    `json:"timeout"`
	Method     string `json:"method"`
	Password   string `json:"password"`
}

type ngrokConfig struct {
	AuthToken string `json:"auth_token"`
}

// LambdaShadowSocksConfig represents the configuration needed for lambda
type LambdaShadowSocksConfig struct {
	Addr      string `json:"addr"`
	Method    string `json:"method"`
	Password  string `json:"password"`
	ProxyHost string `json:"proxyHost"`
	ProxyPort string `json:"proxyPort"`
}

// Config represents the project's configuration
type Config struct {
	AWS         *awsConfig         `json:"AWS"`
	Shadowsocks *shadowsocksConfig `json:"shadowsocks"`
	Ngrok       *ngrokConfig       `json:"ngrok"`
}

// FuncConfig represents the slss lambda function configuration
type FuncConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Runtime     string `json:"runtime"`
	Memory      int    `json:"memory"`
	Timeout     int    `json:"timeout"`
}

// LoadFuncConfig loads the lambda function configuration from a specified path
func LoadFuncConfig(path string) (*FuncConfig, error) {
	var config = new(FuncConfig)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, errors.WithStack(err)
	}

	if config.Timeout < 60 {
		return nil, errors.New("timeout in function configuration should >= 60")
	}

	return config, nil
}

// LoadConfig loads the configuration object from a specified path
func LoadConfig(path string) (*Config, error) {
	var config = new(Config)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, errors.WithStack(err)
	}

	// Try to find AWS configuration from environment variables
	if config.AWS.AccessKeyID == "" {
		config.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if config.AWS.SecretAccessKey == "" {
		config.AWS.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	if config.AWS.Region == "" {
		config.AWS.Region = os.Getenv("AWS_REGION")
	}

	if config.AWS.AccessKeyID == "" || config.AWS.SecretAccessKey == "" || config.AWS.Region == "" {
		return nil, errors.New("empty AWS configuration")
	}

	return config, nil
}
