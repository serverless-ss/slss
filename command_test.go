package slss

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommandSuite struct {
	suite.Suite
	APEXExecutor *APEXCommandExecutor
}

func (s *CommandSuite) SetupTest() {
	s.APEXExecutor = &APEXCommandExecutor{Config: &Config{
		AWS: &awsConfig{
			AccessKeyID:     "KEY_ID",
			SecretAccessKey: "SECRET",
			Region:          "REGION",
		},
	}}
}

func (s *CommandSuite) TestAPEXCommandExecutorExec() {
	output, err := s.APEXExecutor.Exec("apex", nil)
	s.Nil(err)
	s.Contains(output, "apex [command]")
}

func TestCommand(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}
