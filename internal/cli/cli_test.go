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

	err := f()

	// Restore original stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf strings.Builder
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestRunIndexCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test input file
	inputPath := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("This is a test file content for indexing"), 0o644); err != nil {
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
				return os.MkdirAll(indexDir, 0o755)
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

	if err := os.WriteFile(input1Path, []byte("This is the first test file"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(input2Path, []byte("This is the second test file"), 0o644); err != nil {
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
	if err := os.WriteFile(inputPath, []byte("This is test content"), 0o644); err != nil {
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

func TestRunFuzzyCommand(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath, _ := createValidIndex(t, tmpDir)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name: "fuzzy search with valid hash",
			args: []string{
				"program",
				"-c", "fuzzy",
				"-i", inputPath,
				"-h", "123456789abcdef0",
				"-threshold", "3",
			},
			wantErr: false,
		},
		{
			name: "fuzzy search with invalid hash format",
			args: []string{
				"program",
				"-c", "fuzzy",
				"-i", inputPath,
				"-h", "invalid",
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

func TestRunModerateCommand(t *testing.T) {
	tmpDir := t.TempDir()

	inputPath := filepath.Join(tmpDir, "content.txt")
	wordlistPath := filepath.Join(tmpDir, "wordlist.txt")

	if err := os.WriteFile(inputPath, []byte("This is test content with forbidden words"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wordlistPath, []byte("forbidden\nrestricted"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name: "moderate with valid inputs",
			args: []string{
				"program",
				"-c", "moderate",
				"-i", inputPath,
				"-wordlist", wordlistPath,
				"-level", "strict",
				"-context", "50",
			},
			wantErr: true,
		},
		{
			name: "moderate without wordlist",
			args: []string{
				"program",
				"-c", "moderate",
				"-i", inputPath,
			},
			wantErr: true,
		},
		{
			name: "moderate with invalid level",
			args: []string{
				"program",
				"-c", "moderate",
				"-i", inputPath,
				"-wordlist", wordlistPath,
				"-level", "invalid_level",
			},
			wantErr: true,
		},
		{
			name: "moderate with nonexistent input",
			args: []string{
				"program",
				"-c", "moderate",
				"-i", "nonexistent.txt",
				"-wordlist", wordlistPath,
				"-level", "strict",
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

func TestRunStatsCommand(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a test index and get both the path and hash
	inputPath, _ := createValidIndex(t, tmpDir)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		setup   func() error
		verify  func() error
	}{
		{
			name: "stats with valid index file",
			args: []string{
				"program",
				"-c", "stats",
				"-i", inputPath,
				"-v",
			},
			wantErr: false,
		},
		{
			name: "stats with logging enabled",
			args: []string{
				"program",
				"-c", "stats",
				"-i", inputPath,
				"-log", filepath.Join(tmpDir, "stats.log"),
			},
			wantErr: false,
			verify: func() error {
				_, err := os.Stat(filepath.Join(tmpDir, "stats.log"))
				return err
			},
		},
		{
			name: "stats with invalid log path",
			args: []string{
				"program",
				"-c", "stats",
				"-i", inputPath,
				"-log", "/invalid/path/stats.log",
			},
			wantErr: true,
		},
		{
			name: "stats with corrupted index file",
			args: []string{
				"program",
				"-c", "stats",
				"-i", filepath.Join(tmpDir, "corrupted.idx"),
			},
			wantErr: true,
			setup: func() error {
				return os.WriteFile(filepath.Join(tmpDir, "corrupted.idx"), []byte("corrupted data"), 0o644)
			},
		},
		{
			name: "stats with directory instead of file",
			args: []string{
				"program",
				"-c", "stats",
				"-i", tmpDir,
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

func TestRunLookupCommand(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath, validHash := createValidIndex(t, tmpDir)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name: "lookup with valid hash",
			args: []string{
				"program",
				"-c", "lookup",
				"-i", inputPath,
				"-h", validHash,
			},
			wantErr: false,
		},
		// Other test cases...
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

func createValidIndex(t *testing.T, tmpDir string) (string, string) {
	inputFile := filepath.Join(tmpDir, "sample.txt")
	if err := os.WriteFile(inputFile, []byte("Sample content for indexing"), 0o644); err != nil {
		t.Fatal("Failed to create sample input file:", err)
	}
	indexFile := filepath.Join(tmpDir, "index.idx")
	err := Run([]string{"program", "-c", "index", "-i", inputFile, "-o", indexFile})
	if err != nil {
		t.Fatal("Failed to create index:", err)
	}
	// Calculate the hash of the input file
	output, err := captureOutput(func() error {
		return Run([]string{"program", "-c", "hash", "-i", inputFile})
	})
	if err != nil {
		t.Fatal("Failed to calculate hash:", err)
	}
	hash := strings.TrimSpace(output)
	if hash == "" {
		t.Fatal("Failed to obtain hash from hash command")
	}
	return indexFile, hash
}
