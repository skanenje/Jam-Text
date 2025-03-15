package main
// Provides the entry point for the JamText application

import (
	"fmt"
	"os"

	"jamtext/internal/cli"
)

func main() {
	// RUN the CLI application
	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
