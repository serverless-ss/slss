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
func (a *APEXCommandExecutor) Exec(command string, stdin *bytes.Buffer, args ...string) (string, error) {
	var (
		responseMessage bytes.Buffer
		errorMessage    bytes.Buffer
	)

	cmd := exec.Command(command, args...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	cmd.Stdout = &responseMessage
	cmd.Stderr = &errorMessage

	wd, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}

	cmd.Dir = wd + "/lambda/"
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf(awsAccessKeyIDTemplate, a.Config.AWS.AccessKeyID),
		fmt.Sprintf(awsSecretAccessKeyTemplate, a.Config.AWS.SecretAccessKey),
		fmt.Sprintf(awsRegionTemplate, a.Config.AWS.Region),
	)

	if err := cmd.Start(); err != nil {
		return "", errors.Wrapf(err, "APEX commend failed: \n%v", errorMessage.String())
	}

	if err := cmd.Wait(); err != nil {
		return "", errors.Wrapf(err, "APEX commend failed: \n%v", errorMessage.String())
	}

	return responseMessage.String(), nil
}
