package chunk

import (
	"testing"
)

func TestDefaultChunkOptions(t *testing.T) {
	opts := DefaultChunkOptions()

	if opts.ChunkSize != 4096 {
		t.Errorf("Expected ChunkSize 4096, got %d", opts.ChunkSize)
	}
	if opts.OverlapSize != 256 {
		t.Errorf("Expected OverlapSie 256, got %d", opts.OverlapSize)
	}

	if !opts.SplitOnBoundary {
		t.Errorf("Expected SplitOnBoundary to be true")
	}
}
