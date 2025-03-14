package main

import (
	"fmt"
	"os"

	"jamtext/internal/cli"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
