package cli

import (
	"flag"
	"fmt"
)

func Run(args []string) error {
	fs := flag.NewFlagSet("textindex", flag.ExitOnError)
	
	// Basic commands
	cmd := fs.String("cmd", "", "Command to run")
	input := fs.String("i", "", "Input file path")
	output := fs.String("o", "", "Output file path")

	fs.Parse(args[1:])

	switch *cmd {
	case "index":
		if *input == "" || *output == "" {
			return fmt.Errorf("input and output file paths must be specified")
		}

		// TODO: Setup chunk options

	case "lookup":
		// TODO: Add a function to lookup chunks

	default:
		// TODO: Setup chunk options

		
		return fmt.Errorf("unknown command: %s", *cmd)
	}

	return nil
}