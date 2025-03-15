package chunk

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"jamtext/internal/simhash"
)

// TestNewChunkProcessor tests the creation of chunk processors with different worker counts
func TestNewChunkProcessor(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		wantCPUs   bool
	}{
		{
			name:       "zero workers defaults to CPU count",
			numWorkers: 0,
			wantCPUs:   true,
		},
		{
			name:       "positive workers respected",
			numWorkers: 4,
			wantCPUs:   false,
		},
		{
			name:       "negative workers defaults to CPU count",
			numWorkers: -1,
			wantCPUs:   true,
		},
	}

	hyperplanes := simhash.GenerateHyperplanes(128, 64)
	cpuCount := runtime.NumCPU()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := NewChunkProcessor(tt.numWorkers, hyperplanes)
			if cp == nil {
				t.Fatal("Expected non-nil ChunkProcessor")
			}
			defer cp.Close()

			expectedWorkers := tt.numWorkers
			if tt.wantCPUs {
				expectedWorkers = cpuCount
			}

			// Check if the worker pool was created with the correct number of workers
			if cp.pool.workers != expectedWorkers {
				t.Errorf("got %d workers, want %d", cp.pool.workers, expectedWorkers)
			}

			// Verify the result channel has the expected buffer size
			if cap(cp.resultChan) != expectedWorkers*2 {
				t.Errorf("result channel capacity = %d, want %d", cap(cp.resultChan), expectedWorkers*2)
			}
		})
	}
}

func TestProcessChunk(t *testing.T) {
	hyperplanes := simhash.GenerateHyperplanes(128, 64)
	cp := NewChunkProcessor(2, hyperplanes)
	defer cp.Close()

	tests := []struct {
		name    string
		chunk   Chunk
		wantPos int64
	}{
		{
			name: "simple chunk",
			chunk: Chunk{
				Content:     "test content",
				StartOffset: 0,
			},
			wantPos: 0,
		},
		{
			name: "chunk with offset",
			chunk: Chunk{
				Content:     "test content",
				StartOffset: 100,
			},
			wantPos: 100,
		},
		{
			name: "empty chunk",
			chunk: Chunk{
				Content:     "",
				StartOffset: 0,
			},
			wantPos: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp.ProcessChunk(tt.chunk)
			result := <-cp.Results()

			if result.Error != nil {
				t.Errorf("Unexpected error: %v", result.Error)
			}
			if result.Pos != tt.wantPos {
				t.Errorf("Expected position %d, got %d", tt.wantPos, result.Pos)
			}
			if result.Hash == 0 {
				t.Error("Expected non-zero hash")
			}
		})
	}
}

func TestProcessorClose(t *testing.T) {
	hyperplanes := simhash.GenerateHyperplanes(128, 64)
	cp := NewChunkProcessor(2, hyperplanes)

	// Submit some work
	cp.ProcessChunk(Chunk{Content: "test", StartOffset: 0})

	// Consume the result
	result := <-cp.Results()
	if result.Error != nil {
		t.Errorf("Unexpected error: %v", result.Error)
	}

	// Close the processor
	cp.Close()

	// Verify channel is closed
	_, ok := <-cp.Results()
	if ok {
		t.Error("Results channel should be closed")
	}
}

func TestProcessorWithEmptyHyperplanes(t *testing.T) {
	cp := NewChunkProcessor(2, nil)
	defer cp.Close()

	cp.ProcessChunk(Chunk{Content: "test", StartOffset: 0})
	result := <-cp.Results()

	if result.Error != nil {
		t.Errorf("Expected no error with empty hyperplanes, got %v", result.Error)
	}
}

func TestProcessorConcurrency(t *testing.T) {
	hyperplanes := simhash.GenerateHyperplanes(128, 64)
	cp := NewChunkProcessor(4, hyperplanes)
	defer cp.Close()

	numChunks := 100
	processed := make(chan struct{}, numChunks)

	// Submit chunks concurrently
	for i := 0; i < numChunks; i++ {
		go func(i int) {
			cp.ProcessChunk(Chunk{
				Content:     fmt.Sprintf("chunk-%d", i),
				StartOffset: int64(i * 100),
			})
		}(i)
	}

	// Collect results
	for i := 0; i < numChunks; i++ {
		result := <-cp.Results()
		if result.Error != nil {
			t.Errorf("Unexpected error processing chunk: %v", result.Error)
		}
		processed <- struct{}{}
	}

	if len(processed) != numChunks {
		t.Errorf("Expected %d processed chunks, got %d", numChunks, len(processed))
	}
}

func TestProcessorWithInvalidContent(t *testing.T) {
	hyperplanes := simhash.GenerateHyperplanes(128, 64)
	cp := NewChunkProcessor(2, hyperplanes)
	defer cp.Close()

	tests := []struct {
		name    string
		content string
	}{
		{"empty content", ""},
		{"very large content", strings.Repeat("a", 1000000)},
		{"special characters", "⌘\n\t\r"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp.ProcessChunk(Chunk{Content: tt.content, StartOffset: 0})
			result := <-cp.Results()

			if result.Error != nil {
				t.Errorf("Unexpected error: %v", result.Error)
			}
			if result.Hash == 0 {
				t.Error("Expected non-zero hash even for edge cases")
			}
		})
	}
}
