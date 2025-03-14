package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestMainFunction(t *testing.T) {
	// Create temp test directory
	tmpDir := t.TempDir()

	// Create test files
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "output.idx")

	// Write test content
	testContent := []byte("This is a test document for indexing.")
	if err := os.WriteFile(inputFile, testContent, 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Save original args and stdout
	oldArgs := os.Args
	oldStdout := os.Stdout
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	// Set up command line args for index command
	os.Args = []string{
		"./jamtext",
		"index",
		"-i", inputFile,
		"-o", outputFile,
	}

	// Capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run main
	main()

	// Restore stdout and get output
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Verify no usage message
	if bytes.Contains(buf.Bytes(), []byte("Usage:")) {
		t.Error("Should not print usage when valid arguments provided")
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Save original args and stderr
	originalArgs := os.Args
	originalStderr := os.Stderr

	// Restore original values when test completes
	defer func() {
		os.Args = originalArgs
		os.Stderr = originalStderr
	}()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no command",
			args:    []string{"./jamtext"},
			wantErr: true,
		},
		{
			name:    "invalid command",
			args:    []string{"./jamtext", "-c", "invalid"},
			wantErr: true,
		},
		{
			name: "valid index command",
			args: []string{
				"./jamtext",
				"-c", "index",
				"-i", inputFile,
				"-o", outputFile,
				"-s", "1024",
			},
			wantErr: false,
		},
		{
			name: "index command missing input",
			args: []string{
				"./jamtext",
				"-c", "index",
				"-o", outputFile,
			},
			wantErr: true,
		},
		{
			name: "hash command",
			args: []string{
				"./jamtext",
				"-c", "hash",
				"-i", inputFile,
			},
			wantErr: false,
		},
		{
			name: "stats command missing input",
			args: []string{
				"./jamtext",
				"-c", "stats",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create pipe to capture stderr
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stderr = w

			// Set test args
			os.Args = tt.args

			// Instead of goroutine, capture exit code
			var buf bytes.Buffer
			exitCode := 0

			// Redirect stdout too if needed
			oldStdout := os.Stdout
			stdoutR, stdoutW, _ := os.Pipe()
			os.Stdout = stdoutW

			func() {
				defer func() {
					if r := recover(); r != nil {
						exitCode = 1
					}
					w.Close()
					stdoutW.Close()
					os.Stderr = originalStderr
					os.Stdout = oldStdout
				}()
				main()
			}()

			// Read outputs
			buf.ReadFrom(r)
			var stdoutBuf bytes.Buffer
			stdoutBuf.ReadFrom(stdoutR)

			// Check if we got an error (exit code != 0) when we wanted one
			gotErr := exitCode != 0
			if gotErr != tt.wantErr {
				t.Errorf("main() error = %v, wantErr %v\nStderr: %s\nStdout: %s",
					gotErr, tt.wantErr, buf.String(), stdoutBuf.String())
			}
		})
	}
}
