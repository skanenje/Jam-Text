package index

import (
	"path/filepath"
	"testing"

	"jamtext/internal/simhash"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := "test.txt"
	chunkSize := 4096
	hyperplanes := simhash.GenerateHyperplanes(128, 64)

	tests := []struct {
		name      string
		indexDir  string
		wantEmpty bool
	}{
		{
			name:      "with custom index directory",
			indexDir:  filepath.Join(tmpDir, "custom"),
			wantEmpty: false,
		},
		{
			name:     "with custom index directory",
			indexDir: filepath.Join(tmpDir, "custome"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := New(sourceFile, chunkSize, hyperplanes, tt.indexDir)

			if idx == nil {
				t.Fatal("Expected non-nil Index")
			}

			if idx.SourceFile != sourceFile {
				t.Errorf("Expected source file %s, got %s", sourceFile, idx.SourceFile)
			}

			if idx.ChunkSize != chunkSize {
				t.Errorf("Expected chunk size %d, got %d", chunkSize, idx.ChunkSize)
			}

			if len(idx.Shards) != 1 {
				t.Errorf("Expected 1 initial shard, got %d", len(idx.Shards))
			}

			if idx.Shards[0] == nil || len(idx.Shards[0].SimHashToPos) != 0 {
				t.Errorf("Expected empty initial shard")
			}
		})
	}
}
