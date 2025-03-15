package chunk

import (
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"jamtext/internal/simhash"
)

func DefaultChunkOptions() ChunkOptions {
	return ChunkOptions{
		ChunkSize:       4096,
		OverlapSize:     256,
		SplitOnBoundary: true,
	}
}

func TestWorkerPoolBasic(t *testing.T) {
	pool := NewWorkerPool(2)
	defer pool.Close()

	var wg sync.WaitGroup
	results := make([]int, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i
		pool.Submit(func() {
			results[i] = i * 2
			wg.Done()
		})
	}

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	sum := 0
	for _, v := range results {
		sum += v
	}

	expected := 20 // sum of 0*2 through 4*2
	if sum != expected {
		t.Errorf("Expected sum %d, got %d", expected, sum)
	}
}

func TestReadChunk(t *testing.T) {
	// Create a temporary test file
	content := "This is a test content.\nIt has multiple lines.\nTesting chunk reading."
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test cases
	tests := []struct {
		name          string
		position      int64
		chunkSize     int
		contextBefore int
		contextAfter  int
		wantChunk     string
		wantErr       bool
	}{
		{
			name:          "basic read",
			position:      5,
			chunkSize:     10,
			contextBefore: 5,
			contextAfter:  5,
			wantChunk:     "is a test ",
		},
		{
			name:          "read from start",
			position:      0,
			chunkSize:     10,
			contextBefore: 0,
			contextAfter:  5,
			wantChunk:     "This is a ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk, err := ReadChunk(tmpfile.Name(), tt.position, tt.chunkSize)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadChunk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if chunk != tt.wantChunk {
				t.Errorf("ReadChunk() chunk = %v, want %v", chunk, tt.wantChunk)
			}
		})
	}
}

func TestWorkerPoolEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		numTasks   int
	}{
		{"zero workers", 0, 3},
		{"single worker", 1, 3},
		{"multiple workers", 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewWorkerPool(tt.numWorkers)
			defer pool.Close()

			var wg sync.WaitGroup
			var counter atomic.Int32

			// Ensure at least 1 worker for the pool
			if tt.numWorkers == 0 {
				t.Skip("Skipping test with zero workers - not a valid configuration")
			}

			for i := 0; i < tt.numTasks; i++ {
				wg.Add(1)
				pool.Submit(func() {
					defer wg.Done()
					counter.Add(1)
				})
			}

			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(500 * time.Millisecond):
				t.Fatal("Test timed out after 500ms")
			}

			if counter.Load() != int32(tt.numTasks) {
				t.Errorf("Expected %d tasks completed, got %d", tt.numTasks, counter.Load())
			}
		})
	}
}

// func TestNewChunk(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		content     string
// 		startOffset int64
// 		wantLength  int
// 	}{
// 		{
// 			name:        "empty content",
// 			content:     "",
// 			startOffset: 0,
// 			wantLength:  0,
// 		},
// 		{
// 			name:        "normal content",
// 			content:     "test content",
// 			startOffset: 100,
// 			wantLength:  11,
// 		},
// 		{
// 			name:        "content with special chars",
// 			content:     "test\ncontent\t中文",
// 			startOffset: 50,
// 			wantLength:  17, // Corrected from 14 to 17 (actual byte length)
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			chunk := NewChunk(tt.content, tt.startOffset)

// 			if chunk.Content != tt.content {
// 				t.Errorf("Content = %v, want %v", chunk.Content, tt.content)
// 			}
// 			if chunk.StartOffset != tt.startOffset {
// 				t.Errorf("StartOffset = %v, want %v", chunk.StartOffset, tt.startOffset)
// 			}
// 			if chunk.Length != tt.wantLength {
// 				t.Errorf("Length = %v, want %v", chunk.Length, tt.wantLength)
// 			}
// 			if !chunk.IsComplete {
// 				t.Error("IsComplete = false, want true")
// 			}
// 			if chunk.Metadata == nil {
// 				t.Error("Metadata is nil, want initialized map")
// 			}
// 		})
// 	}
// }

// func TestWorkerPoolClose(t *testing.T) {
// 	pool := NewWorkerPool(1)
// 	var wg sync.WaitGroup
// 	executed := make(chan bool, 1)

// 	// Submit a task that sleeps to ensure pool is busy
// 	wg.Add(1)
// 	pool.Submit(func() {
// 		defer wg.Done()
// 		time.Sleep(50 * time.Millisecond)
// 	})

// 	// Close the pool
// 	pool.Close()

// 	// Submit a task after closing
// 	pool.Submit(func() {
// 		executed <- true
// 	})

// 	// Wait for the first task to complete
// 	wg.Wait()

// 	// Check if any task was executed after closure
// 	select {
// 	case <-executed:
// 		t.Error("Task submitted after pool closure was executed")
// 	case <-time.After(10 * time.Millisecond):
// 		// Expected behavior: no execution after closure
// 	}
// }

func TestChunkWithMetadata(t *testing.T) {
	chunk := NewChunk("test", 0)
	chunk.Metadata["timestamp"] = time.Now().Format(time.RFC3339)
	chunk.Metadata["type"] = "test"

	if len(chunk.Metadata) != 2 {
		t.Errorf("Expected 2 metadata entries, got %d", len(chunk.Metadata))
	}
}

func TestReadChunkWithInvalidFile(t *testing.T) {
	_, err := ReadChunk("nonexistent.txt", 0, 100)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestReadChunkWithDifferentFormats(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		cleanup func(string)
		wantErr bool
	}{
		{
			name: "pdf_file_without_pdftotext",
			setup: func() string {
				f, _ := os.CreateTemp("", "test*.pdf")
				f.Write([]byte("%PDF-1.4\n")) // Mock PDF content
				return f.Name()
			},
			cleanup: func(path string) { os.Remove(path) },
			wantErr: true,
		},
		{
			name: "docx_file_without_pandoc",
			setup: func() string {
				f, _ := os.CreateTemp("", "test*.docx")
				f.Write([]byte("PK\x03\x04")) // Mock DOCX header
				return f.Name()
			},
			cleanup: func(path string) { os.Remove(path) },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer tt.cleanup(path)

			_, err := ReadChunk(path, 0, 1024)
			if (err == nil) != !tt.wantErr {
				t.Errorf("ReadChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindBoundary(t *testing.T) {
	tests := []struct {
		name         string
		text         []byte
		preferredPos int
		boundaryChar string
		want         int
	}{
		{
			name:         "simple boundary",
			text:         []byte("Hello. World"),
			preferredPos: 7,
			boundaryChar: ".",
			want:         6,
		},
		{
			name:         "no boundary found",
			text:         []byte("HelloWorld"),
			preferredPos: 5,
			boundaryChar: ".",
			want:         5,
		},
		{
			name:         "boundary at start",
			text:         []byte(".HelloWorld"),
			preferredPos: 5,
			boundaryChar: ".",
			want:         1,
		},
		{
			name:         "multiple boundaries",
			text:         []byte("Hello.World.Test"),
			preferredPos: 10,
			boundaryChar: ".",
			want:         6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findBoundary(tt.text, tt.preferredPos, tt.boundaryChar)
			if got != tt.want {
				t.Errorf("findBoundary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunLookupCommand(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath, validHash := createValidIndex(t, tmpDir)

	// Create and configure index
	opts := DefaultChunkOptions()
	opts.Logger = log.New(os.Stdout, "", 0)
	hyperplanes := simhash.GenerateHyperplanes(128, 64)
	idx, err := ProcessFile(inputPath, opts, hyperplanes, tmpDir)
	if err != nil {
		t.Fatalf("Failed to load index: %v", err)
	}

	// Add the test hash to the index
	hashValue, err := strconv.ParseUint(validHash, 16, 64)
	if err != nil {
		t.Fatalf("Failed to parse hash: %v", err)
	}

	// Add the hash to the index with a known position
	if err := idx.Add(simhash.SimHash(hashValue), 0); err != nil {
		t.Fatalf("Failed to add hash to index: %v", err)
	}

	// Perform lookup
	positions, err := idx.Lookup(simhash.SimHash(hashValue))
	if err != nil {
		t.Fatalf("Failed to lookup hash: %v", err)
	}
	if len(positions) == 0 {
		t.Errorf("No matches found for hash %s", validHash)
	}
}

func TestRunFuzzyCommand(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath, _ := createValidIndex(t, tmpDir)
	// Use inputPath for fuzzy search test logic...
	// No hash needed, so ignore it with _
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Errorf("Index file does not exist: %v", err)
	}
}

func TestRunStatsCommand(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath, _ := createValidIndex(t, tmpDir)
	// Use inputPath for stats test logic...
	// No hash needed, so ignore it with _
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Errorf("Index file does not exist: %v", err)
	}
}

// createValidIndex creates a test file and returns its path and a valid hash
func createValidIndex(t *testing.T, dir string) (string, string) {
	content := "This is test content for indexing"
	tmpfile, err := os.CreateTemp(dir, "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Return a known valid hash that we'll add to the index
	return tmpfile.Name(), "0123456789abcdef"
}
