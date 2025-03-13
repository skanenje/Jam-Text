package main

import (
	cli "jamtext/internal/cli"
	"os"
)

func main() {
	// TODO
	if err := cli.Run(os.Args); err != nil {
		panic(err)
	}
}