package slss

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

const projectConfigPath = "./lambda/project.json"

// Amazon AWS configuration
type awsConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Region          string `json:"region"`
	Role            string `json:"role"`
}

// Shadowsocks configuration
type shadowsocksConfig struct {
	LocalPort string `json:"local_port"`
	Timeout   int    `json:"timeout"`
	Method    string `json:"method"`
	Password  string `json:"password"`
}

type ngrokConfig struct {
	AuthToken string `json:"auth_token"`
}

// LambdaShadowSocksConfig represents the configuration needed for lambda
type LambdaShadowSocksConfig struct {
	Port       string `json:"port"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	ProxyURL   string `json:"proxyURL"`
	NgrokToken string `json:"ngrokToken"`
}

// Config represents the project's configuration
type Config struct {
	AWS             *awsConfig         `json:"AWS"`
	Shadowsocks     *shadowsocksConfig `json:"shadowsocks"`
	Ngrok           *ngrokConfig       `json:"ngrok"`
	LocalServerPort string             `json:"local_server_port"`
}

// FuncConfig represents the slss lambda function configuration
type FuncConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Runtime     string `json:"runtime"`
	Memory      int    `json:"memory"`
	Timeout     int    `json:"timeout"`
}

// ProjectConfig represents the apex project configuration
type ProjectConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Memory      int    `json:"memory"`
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

// UpdateProjectConfigRole updates the "role" filed in apex project
// configuration
func UpdateProjectConfigRole(role string) error {
	var config = new(ProjectConfig)

	content, err := ioutil.ReadFile(projectConfigPath)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := json.Unmarshal(content, config); err != nil {
		return errors.WithStack(err)
	}

	config.Role = role
	contentToUpdate, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	return ioutil.WriteFile(projectConfigPath, []byte(contentToUpdate), 0644)
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

	if config.AWS == nil || config.Ngrok == nil || config.Shadowsocks == nil {
		return nil, errors.New("empty configuration")
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

	if config.AWS.AccessKeyID == "" ||
		config.AWS.SecretAccessKey == "" ||
		config.AWS.Region == "" ||
		config.AWS.Role == "" {
		return nil, errors.New("empty AWS configuration")
	}

	if config.Ngrok.AuthToken == "" {
		return nil, errors.New("empty ngrok configuration")
	}

	if config.LocalServerPort == "" {
		return nil, errors.New("empty local_server_port configuration")
	}

	return config, nil
}
