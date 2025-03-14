package main

import (
	"bytes"
	"fmt"
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
	
	err := os.WriteFile(inputFile, []byte("This is a test document for indexing."), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
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
		name     string
		args     []string
		wantErr  bool
	}{
		{
			name:    "no command",
			args:    []string{"jamtext"},
			wantErr: true,
		},
		{
			name:    "invalid command",
			args:    []string{"jamtext", "-c", "invalid"},
			wantErr: true,
		},
		{
			name: "valid index command",
			args: []string{
				"jamtext",
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
				"jamtext",
				"-c", "index",
				"-o", outputFile,
			},
			wantErr: true,
		},
		{
			name: "hash command",
			args: []string{
				"jamtext",
				"-c", "hash",
				"-i", inputFile,
			},
			wantErr: false,
		},
		{
			name: "stats command missing input",
			args: []string{
				"jamtext",
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
			
			// Run main and capture error
			var runErr error
			done := make(chan bool)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						runErr = fmt.Errorf("panic: %v", r)
					}
					done <- true
				}()
				main()
			}()

			// Wait for completion
			<-done
			
			// Close pipe and read output
			w.Close()
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("Failed to read pipe: %v", err)
			}
			
			// Check error status
			if (runErr != nil) != tt.wantErr {
				t.Errorf("main() error = %v, wantErr %v\nOutput: %s", 
					runErr, tt.wantErr, buf.String())
			}
		})
	}
}
