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
			args:    []string{"program", "-c", "invalid"},
			wantErr: true,
			errMsg:  "unknown command: invalid",
		},
		{
			name:    "index without input",
			args:    []string{"program", "-c", "index"},
			wantErr: true,
			errMsg:  "input and output file paths must be specified",
		},
		{
			name:    "index without output",
			args:    []string{"program", "-c", "index"},
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
	if err := os.WriteFile(inputPath, []byte("This is a test file content for indexing"), 0o644); err != nil {
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
				"-c", "index",
				"-i", inputPath,
				"-o", outputPath,
			},
			wantErr: false,
		},
		{
			name: "index with all options",
			args: []string{
				"program",
				"-c", "index",
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
				return os.MkdirAll(indexDir, 0o755)
			},
			cleanup: func() error {
				return os.RemoveAll(indexDir)
			},
		},
		{
			name: "index with invalid input file",
			args: []string{
				"program",
				"-c", "index",
				"-i", "non existent.txt",
				"-o", outputPath,
			},
			wantErr: true,
		},
		{
			name: "index with invalid log file path",
			args: []string{
				"program",
				"-c", "index",
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
	if err := os.WriteFile(inputPath, []byte("This is test content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create index file using the index command
	indexPath := filepath.Join(tmpDir, "test.idx")
	indexArgs := []string{
		"program",
		"-c", "index",
		"-i", inputPath,
		"-o", indexPath,
	}
	if err := Run(indexArgs); err != nil {
		t.Fatal(err)
	}

	logPath := filepath.Join(tmpDir, "lookup.log")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name: "successful lookup with defaults",
			args: []string{
				"program",
				"-c", "lookup",
				"-i", indexPath,
				"-h", "108fb9408bf49bee",
			},
			wantErr: false,
		},
		{
			name: "lookup with all options",
			args: []string{
				"program",
				"-c", "lookup",
				"-i", indexPath,
				"-h", "108fb9408bf49bee",
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
				"-c", "lookup",
				"-i", "nonexistent.idx",
				"-h", "108fb9408bf49bee",
			},
			wantErr: true,
		},
		{
			name: "lookup without hash value",
			args: []string{
				"program",
				"-c", "lookup",
				"-i", indexPath,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Clean up
	os.RemoveAll(tmpDir)
}
