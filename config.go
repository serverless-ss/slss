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
	Region          string `json:"Region"`
}

// Config represents the project's configuration
type Config struct {
	AWS awsConfig `json:"AWS"`
}

// LoadConfig loads the configuration object from a specified path
func LoadConfig(path string) (*Config, error) {
	var config = new(Config)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read configuration file failed")
	}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, errors.Wrap(err, "unmarshal configuration file's content failed")
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
