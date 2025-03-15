package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	// Build the executable
	cmd := exec.Command("go", "build", "-o", "jamtext")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to build executable: %v", err)
	}

	// Run the tests
	code := m.Run()
	os.Remove("jamtext")
	os.Exit(code)
}

func TestMainFunction(t *testing.T) {
	// Build the executable first
	buildCmd := exec.Command("go", "build", "-o", "jamtext")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build executable: %v", err)
	}
	defer os.Remove("jamtext")

	// Setup: create a temporary input file
	inputFile := "input.txt"
	outputFile := "output.txt"
	if err := os.WriteFile(inputFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}
	defer os.Remove(inputFile)
	defer os.Remove(outputFile) // Clean up output file if created

	// Define test cases
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		checkOutput func(output []byte) error
	}{
		{
			name:    "valid index command",
			args:    []string{"./jamtext", "-c", "index", "-i", inputFile, "-o", outputFile},
			wantErr: false,
			checkOutput: func(output []byte) error {
				if bytes.Contains(output, []byte("Usage:")) {
					return fmt.Errorf("should not print usage for valid command")
				}
				if _, err := os.Stat(outputFile); os.IsNotExist(err) {
					return fmt.Errorf("output file was not created")
				}
				return nil
			},
		},
		{
			name:    "no_command_provided",
			args:    []string{"./jamtext"},
			wantErr: false, // Changed from true to false
			checkOutput: func(output []byte) error {
				if !bytes.Contains(output, []byte("Usage:")) {
					return fmt.Errorf("should print usage when no command is provided")
				}
				return nil
			},
		},
	}

	// Run table-driven tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(tt.args[0], tt.args[1:]...)
			output, err := cmd.CombinedOutput()
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("Command %v: got error %v, want error %v\nOutput: %s", tt.args, err, tt.wantErr, output)
			}
			if tt.checkOutput != nil {
				if checkErr := tt.checkOutput(output); checkErr != nil {
					t.Errorf("Output check failed: %v\nOutput: %s", checkErr, output)
				}
			}
		})
	}
}
