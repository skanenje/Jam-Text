package chunk

import (
	"os"
	"testing"
)


func DefaultChunkOptions() ChunkOptions {
	return ChunkOptions{
		ChunkSize:       4096,
		OverlapSize:     256,
		SplitOnBoundary: true,
	}
}

// func TestDefaultChunkOptions(t *testing.T) {
// 	opts := DefaultChunkOptions()

// 	if opts.ChunkSize != 4096 {
// 		t.Errorf("Expected ChunkSize 4096, got %d", opts.ChunkSize)
// 	}
// 	if opts.OverlapSize != 256 {
// 		t.Errorf("Expected OverlapSie 256, got %d", opts.OverlapSize)
// 	}

// 	if !opts.SplitOnBoundary {
// 		t.Errorf("Expected SplitOnBoundary to be true")
// 	}
// }

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Close()

	results := make(chan int, 10)
	for i := 0; i < 10; i++ {
		i := i
		pool.Submit(func() {
			results <- i * 2
		})
	}

	// Collecting results
	sum := 0
	for i := 0; i < 10; i++ {
		sum += <-results
	}

	expected := 90 // sum of 0*2 through 9*2
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