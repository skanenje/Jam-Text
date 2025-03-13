package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no command",
			args:    []string{"program"},
			wantErr: true,
			errMsg:  "unknown command",
		},
		{
			name:    "ivalid command",
			args:    []string{"program", "-cmd", "invalid"},
			wantErr: true,
			errMsg:  "unknown command: invalid",
		},
		{
			name:    "index without input",
			args:    []string{"program", "-cmd", "index"},
			wantErr: true,
			errMsg:  "input and output file paths must be specified",
		},
		{
			name:    "index without output",
			args:    []string{"program", "-cmd", "index"},
			wantErr: true,
			errMsg:  "input and output file paths must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Run() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
}

func TestRunIndexCommnd(t *testing.T) {
	// Temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "cli_test")
	if err !=  nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test input file
	inputPath := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("This is a test file content for indexing"), 0644); err != nil {
		t.Fatal(err)
	}

	// Creating test log file path
	logPath := filepath.Join(tmpDir, "input.txt")
	outputPath := filepath.Join(tmpDir, "output.idx")
	indexDir := filepath.Join(tmpDir, "index")

	tests := []struct {
		name string
		args []string
		wantErr bool
		setup func() error 
		cleanup func() error
	}{
		{
			name: "successful index with defautls",
			args: []string{
				"program",
				"-cmd", "index",
				"-i", inputPath,
				"-o", outputPath,
			},
			wantErr: false,
		},
		{
			name: "index with all options",
			args: []string{
				"program",
				"-cmd", "index",
				"-i", inputPath,
				"-o", outputPath,
				"-v",
				"-s", "2048",
				"-overlap", "128",
				"-boundary=true",
				"-boundary-chars", ".!?",
				"-max-size", "4096",
				"-preserve-nl=true",
				"-index-dir", indexDir,
			},
			wantErr: false,
			setup: func() error {
				return os.MkdirAll(indexDir, 0755)
			},
			cleanup: func() error {
				return os.RemoveAll(indexDir)
			},
		},
		{
			name: "index with invalid input file",
			args: []string {
				"program",
				"-cmd", "index",
				"-i", "non existent.txt",
				"-o", outputPath,
			},
			wantErr: true,
		},
		{
			name: "index with invalid log file path",
			args: []string{
				"program",
				"-cmd", "index",
				"-i", inputPath,
				"-o", outputPath,
				"-log", "/invalid/path/log.txt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatal(err)
				}
			}

			err := Run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.cleanup != nil {
				if err := tt.cleanup(); err != nil {
					t.Fatal(err)
				}
			}

			// Clean up output files after each test
			os.Remove(outputPath)
			os.Remove(logPath)
		})
	}
}