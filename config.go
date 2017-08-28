package slss

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Amazon AWS configuration
type awsConfig struct {
	accessKeyID     string
	secretAccessKey string
	region          string
}

// Config represents the project's configuration
type Config struct {
	AWS *awsConfig
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
	if config.AWS.accessKeyID == "" {
		config.AWS.accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if config.AWS.secretAccessKey == "" {
		config.AWS.secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	if config.AWS.region == "" {
		config.AWS.region = os.Getenv("AWS_REGION")
	}

	if config.AWS.accessKeyID == "" || config.AWS.secretAccessKey == "" || config.AWS.region == "" {
		return nil, errors.New("empty AWS configuration")
	}

	return config, nil
}
