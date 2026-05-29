// entire-run is an Entire CLI external command.
//
// Once built as an executable named `entire-run`, the parent Entire CLI
// dispatches it when a user runs `entire run`.
package main

import (
	"fmt"
	"os"

	"github.com/entireio/entire-run/internal/cli"
)

var version = "dev"

func main() {
	if err := cli.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
