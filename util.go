package slss

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// PrintErrorAndExit prints the verbose error message and then exit with -1
// error code
func PrintErrorAndExit(err error) {
	log.Errorf("%+v\n", err)
	os.Exit(1)
}
