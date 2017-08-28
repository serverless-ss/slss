package slss

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
}

func (s *UtilSuite) TestPrintErrorAndExit() {
	s.Panics(func() {
		PrintErrorAndExit(errors.New("test-error"))
	})
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}
