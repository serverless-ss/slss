package slss

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

const (
	awsAccessKeyIDTemplate     = "AWS_ACCESS_KEY_ID=%v"
	awsSecretAccessKeyTemplate = "AWS_SECRET_ACCESS_KEY=%v"
	awsRegionTemplate          = "AWS_REGION=%v"
)

// APEXCommandExecutor represents the APEX command executor
type APEXCommandExecutor struct {
	Config *Config
}

// Exec executes the specified APEX command
func (a *APEXCommandExecutor) Exec(command string) (string, error) {
	var responseMessage bytes.Buffer

	cmd := exec.Command(command)
	cmd.Stdout = &responseMessage

	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "get working dir failed")
	}

	cmd.Dir = wd + "/lambda/"
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf(awsAccessKeyIDTemplate, a.Config.AWS.AccessKeyID),
		fmt.Sprintf(awsSecretAccessKeyTemplate, a.Config.AWS.SecretAccessKey),
		fmt.Sprintf(awsRegionTemplate, a.Config.AWS.Region),
	)

	if err := cmd.Run(); err != nil {
		return "", errors.Wrap(err, "APEX command failed")
	}

	return responseMessage.String(), nil
}
