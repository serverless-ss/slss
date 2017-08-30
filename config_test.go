package slss

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (s *ConfigSuite) TestLoadConfig() {
	config, err := LoadConfig("./config.example.json")

	s.Nil(err)
	s.Equal(config.AWS.AccessKeyID, "YOUR_AWS_ACCESS_KEY_ID")
	s.Equal(config.AWS.SecretAccessKey, "YOUR_AWS_SECRET_ACCESS_KEY")
	s.Equal(config.AWS.Region, "YOUR_AWS_REGION")
	s.Equal(config.Shadowsocks.LocalAddr, "LOCAL_ADDR")
	s.Equal(config.Shadowsocks.ServerAddr, "SERVER_ADDR")
	s.Equal(config.Shadowsocks.Timeout, 300)
	s.Equal(config.Shadowsocks.Method, "aes-128-cfb")
	s.Equal(config.Shadowsocks.Password, "PASSWORD")
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
