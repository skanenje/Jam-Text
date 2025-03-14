package chunk

import (
	"runtime"
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
