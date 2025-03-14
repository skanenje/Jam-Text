package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to capture stdout during test execution
func captureOutput(f func() error) (string, error) {
	// Save original stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	err := f()

	// Restore original stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf strings.Builder
	io.Copy(&buf, r)
	return buf.String(), err
}

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
			name:    "invalid command",
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
		// Removed duplicate test case for "index without output"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := captureOutput(func() error {
				return Run(tt.args)
			})
			
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

func TestRunIndexCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test input file
	inputPath := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("This is a test file content for indexing"), 0644); err != nil {
		t.Fatal(err)
	}

	logPath := filepath.Join(tmpDir, "test.log")
	outputPath := filepath.Join(tmpDir, "output.idx")
	indexDir := filepath.Join(tmpDir, "index")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		setup   func() error
		verify  func() error
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
			verify: func() error {
				_, err := os.Stat(outputPath)
				return err
			},
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
				return os.MkdirAll(indexDir, 0755)
			},
			verify: func() error {
				if _, err := os.Stat(outputPath); err != nil {
					return err
				}
				if _, err := os.Stat(logPath); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name: "index with invalid input file",
			args: []string{
				"program",
				"-c", "index",
				"-i", "nonexistent.txt",
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

			_, err := captureOutput(func() error {
				return Run(tt.args)
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && tt.verify != nil {
				if err := tt.verify(); err != nil {
					t.Errorf("verification failed: %v", err)
				}
			}
		})
	}
}

func TestRunCompareCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test input files
	input1Path := filepath.Join(tmpDir, "file1.txt")
	input2Path := filepath.Join(tmpDir, "file2.txt") 
	outputPath := filepath.Join(tmpDir, "report.txt")
	logPath := filepath.Join(tmpDir, "compare.log")

	if err := os.WriteFile(input1Path, []byte("This is the first test file"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(input2Path, []byte("This is the second test file"), 0644); err != nil {
		t.Fatal(err)  
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name: "successful comparison with defaults",
			args: []string{
				"program",
				"-c", "compare", 
				"-i", input1Path,
				"-i2", input2Path,
			},
			wantErr: false,
		},
		{
			name: "comparison with output report",
			args: []string{
				"program", 
				"-c", "compare",
				"-i", input1Path,
				"-i2", input2Path,
				"-o", outputPath,
				"-v",
				"-log", logPath,
			},
			wantErr: false,
		},
		{
			name: "comparison without second input",
			args: []string{
				"program",
				"-c", "compare",
				"-i", input1Path,
			},
			wantErr: true,
		},
		{
			name: "comparison with nonexistent first file",
			args: []string{
				"program",
				"-c", "compare",
				"-i", "nonexistent.txt",
				"-i2", input2Path,
			},
			wantErr: true,
		},
		{
			name: "comparison with nonexistent second file", 
			args: []string{
				"program",
				"-c", "compare",
				"-i", input1Path,
				"-i2", "nonexistent.txt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := captureOutput(func() error {
				return Run(tt.args)
			})
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetContentHash(t *testing.T) {
	tmpDir := t.TempDir()

	inputPath := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("This is test content"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := captureOutput(func() error {
		return Run([]string{
			"program",
			"-c", "hash",
			"-i", inputPath,
		})
	})
	
	if err != nil {
		t.Errorf("GetContentHash() failed: %v", err)
	}
}

