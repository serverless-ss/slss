package slss

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (s *ConfigSuite) TestLoadConfig() {
	config, err := LoadConfig("./config.json")

	s.Nil(err)
	s.Equal(config.AWS.AccessKeyID, "YOUR_AWS_ACCESS_KEY_ID")
	s.Equal(config.AWS.SecretAccessKey, "YOUR_AWS_SECRET_ACCESS_KEY")
	s.Equal(config.AWS.Region, "YOUR_AWS_REGION")
	s.Equal(config.Shadowsocks.LocalPort, "LOCAL_PORT")
	s.Equal(config.Shadowsocks.ServerAddr, "SERVER_ADDR")
	s.Equal(config.Shadowsocks.ServerPort, "SERVER_PORT")
	s.Equal(config.Shadowsocks.Timeout, 300)
	s.Equal(config.Shadowsocks.Method, "aes-128-cfb")
	s.Equal(config.Shadowsocks.Password, "PASSWORD")
}

func (s *ConfigSuite) TestLoadFuncConfig() {
	config, err := LoadFuncConfig("./lambda/functions/slss/function.json")

	s.Nil(err)
	s.Equal(config.Name, "slss")
	s.Equal(config.Description, "slss lambda server function")
	s.Equal(config.Runtime, "nodejs6.10")
	s.Equal(config.Memory, 128)
	s.Equal(config.Timeout, 300)
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
