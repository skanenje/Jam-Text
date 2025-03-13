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
				"-cmd","index",
				"-i", inputPath,
				"-o", outputPath,
				"-v",
				"-s", "2048",
				"-overlap","128",
				"-boundery=true",
				"-boundery-chars", ".!?",
				"-max-size", "4096",
				"-preserve-nl=true",
				"-index-dir", indexDir,
			}
		}
	}
	
}