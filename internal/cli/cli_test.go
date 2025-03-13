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
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test input file
	inputPath := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("This is a test file content for indexing"), 0644); err != nil {
		t.Fatal(err)
	}

	// Creating test log file path
	logPath := filepath.Join(tmpDir, "test.log")
	outputPath := filepath.Join(tmpDir, "output.idx")
	indexDir := filepath.Join(tmpDir, "index")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		setup   func() error
		cleanup func() error
	}{
		{
			name: "successful index with defaults",
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
				"-log", logPath,
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
			args: []string{
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

func TestRunLookupCommand(t *testing.T) {
	// Temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "cli_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a real index file for testing
	inputPath := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("This is test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create index file using the index command
	indexPath := filepath.Join(tmpDir, "test.idx")
	indexArgs := []string{
		"program",
		"-cmd", "index",
		"-i", inputPath,
		"-o", indexPath,
	}
	if err := Run(indexArgs); err != nil {
		t.Fatal(err)
	}

	// Create paths for output and log files
	outputPath := filepath.Join(tmpDir, "lookup_results.txt")
	logPath := filepath.Join(tmpDir, "lookup.log")

	// Use a known test hash - we'll use a dummy value for now
	// In a real scenario, you might want to calculate this based on your SimHash implementation
	testHash := "108fb9408bf49bee" // This is the hash we see in the test output

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		setup   func() error
		cleanup func() error
	}{
		{
			name: "successful lookup with defaults",
			args: []string{
				"program",
				"-cmd", "lookup",
				"-i", indexPath,
				"-h", testHash,
				"-o", outputPath,
			},
			wantErr: false,
		},
		{
			name: "lookup with all options",
			args: []string{
				"program",
				"-cmd", "lookup",
				"-i", indexPath,
				"-h", testHash,
				"-o", outputPath,
				"-v",
				"-log", logPath,
				"-context-before", "200",
				"-context-after", "200",
			},
			wantErr: false,
		},
		{
			name: "lookup with invalid index file",
			args: []string{
				"program",
				"-cmd", "lookup",
				"-i", "nonexistent.idx",
				"-h", testHash,
				"-o", outputPath,
			},
			wantErr: true,
		},
		{
			name: "lookup without hash value",
			args: []string{
				"program",
				"-cmd", "lookup",
				"-i", indexPath,
				"-o", outputPath,
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


