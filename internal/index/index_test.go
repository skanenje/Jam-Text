package index

import (
	"os"
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
			indexDir: filepath.Join(tmpDir, "custom"),
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

func TestAdd(t *testing.T) {
	idx := New("test.txt", 4096, simhash.GenerateHyperplanes(128, 64), "")

	tests := []struct {
		name    string
		hash    simhash.SimHash
		pos     int64
		wantErr bool
	}{
		{
			name:    "add first hash",
			hash:    0x1234,
			pos:     100,
			wantErr: false,
		},
		{
			name:    "add duplicate hash",
			hash:    0x1234,
			pos:     200,
			wantErr: false,
		},
		{
			name:    "add different hash",
			hash:    0x5678,
			pos:     300,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := idx.Add(tt.hash, tt.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			positions, err := idx.Lookup(tt.hash)
			if err != nil {
				t.Fatalf("Lookup failed: %v", err)
			}

			found := false
			for _, pos := range positions {
				if pos == tt.pos {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Postion %d not found for hash %x", tt.pos, tt.hash)
			}
		})
	}
}

func TestSharding(t *testing.T) {
	tmpDir := t.TempDir()
	idx := New("test.txt", 4096, simhash.GenerateHyperplanes(128, 64), tmpDir)

	// Adding MaxShaedSize + 1 entries to force shard rotation
	for i := 0; i < MaxShardSize; i++ {
		hash := simhash.SimHash(i)
		if err := idx.Add(hash, int64(i)); err != nil {
			t.Fatalf("Failed to add hash: %v", err)
		}
	}

	if len(idx.Shards) != 2 {
		t.Errorf("Expected 2 shards after rotation, got %d", len(idx.Shards))
	}

	// Verify first shard was persisted
	shardPath := filepath.Join(tmpDir, idx.ShardFilename+".0")
	if _, err := os.Stat(shardPath); os.IsNotExist(err) {
		t.Errorf("Shard file not created: %s", shardPath)
	}
}
