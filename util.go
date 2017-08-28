package slss

import (
	"fmt"
)

// PrintErrorAndExit prints the verbose error message and then exit with -1
// error code
func PrintErrorAndExit(err error) {
	panic(fmt.Sprintf("%+v\n", err))
}
