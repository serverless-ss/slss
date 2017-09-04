package slss

import (
	"fmt"
	"os"
)

// PrintErrorAndExit prints the verbose error message and then exit with -1
// error code
func PrintErrorAndExit(err error) {
	fmt.Println(fmt.Sprintf("%+v\n", err))
	os.Exit(1)
}
