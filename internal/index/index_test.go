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

func TestFuzzyLookup(t *testing.T) {
	idx := New("test.txt", 4096, simhash.GenerateHyperplanes(128, 64), "")

	// Add test hashes
	testData := []struct {
		hash simhash.SimHash
		pos  int64
	}{
		{0xFF00, 100}, // Base hash
		{0xFF01, 200}, // Distance 1 from 0xFF00
		{0xFF10, 300}, // Distance 1 from 0xFF00
		{0x00FF, 400}, // Different hash
	}

	for _, td := range testData {
		if err := idx.Add(td.hash, td.pos); err != nil {
			t.Fatalf("Failed to add hash: %v", err)
		}
	}

	tests := []struct {
		name          string
		searchHash    simhash.SimHash
		threshold     int
		wantMatches   int
		wantPositions []int64
	}{
		{
			name:          "exact match with threshold 2",
			searchHash:    0xFF00,
			threshold:     2,
			wantMatches:   3,
			wantPositions: []int64{100, 200, 300},
		},
		{
			name:          "no matches",
			searchHash:    0x0000,
			threshold:     1,
			wantMatches:   0,
			wantPositions: nil,
		},
		{
			name:          "single match with threshold 0",
			searchHash:    0xFF01,
			threshold:     0,
			wantMatches:   1,
			wantPositions: []int64{200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, exists := idx.FuzzyLookup(tt.searchHash, tt.threshold)

			if !exists && tt.wantMatches > 0 {
				t.Error("Expected matches but got none")
				return
			}

			totalMatches := 0
			foundPositions := make(map[int64]bool)

			// Count total matches and collect all positions
			for _, positions := range results {
				totalMatches++
				for _, pos := range positions {
					foundPositions[pos] = true
				}
			}

			if totalMatches != tt.wantMatches {
				t.Errorf("Expected %d matches, got %d", tt.wantMatches, totalMatches)
			}

			// Verify each expected position is found
			if tt.wantPositions != nil {
				for _, wantPos := range tt.wantPositions {
					if !foundPositions[wantPos] {
						t.Errorf("Expected position %d not found in results", wantPos)
					}
				}
			}
		})
	}
}

func TestLookupNonexistentHash(t *testing.T) {
	idx := New("test.txt", 4096, simhash.GenerateHyperplanes(128, 64), "")

	tests := []struct {
		name    string
		hash    simhash.SimHash
		wantLen int
	}{
		{
			name:    "lookup nonexistent hash",
			hash:    0xABCD,
			wantLen: 0,
		},
		{
			name:    "lookup after adding different hash",
			hash:    0x1111,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add a different hash first
			if err := idx.Add(0x9999, 100); err != nil {
				t.Fatalf("Failed to add hash: %v", err)
			}

			positions, err := idx.Lookup(tt.hash)
			if err != nil {
				t.Fatalf("Lookup failed: %v", err)
			}

			if len(positions) != tt.wantLen {
				t.Errorf("Expected %d positions, got %d", tt.wantLen, len(positions))
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := "test.txt"
	chunkSize := 4096
	hyperplanes := simhash.GenerateHyperplanes(128, 64)

	// Create and populate index
	idx := New(sourceFile, chunkSize, hyperplanes, tmpDir)
	hash1 := simhash.SimHash(0x1234)
	pos1 := int64(100)
	if err := idx.Add(hash1, pos1); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	hash2 := simhash.SimHash(0x5678)
	pos2 := int64(200)
	if err := idx.Add(hash2, pos2); err != nil {
		t.Fatalf("AAdd faile: %v", err)
	}

	// Save index
	indexFile := filepath.Join(tmpDir, "index.gob")
	if err := Save(idx, indexFile); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	loadedIdx, err := Load(indexFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify loaded data
	positions1, err := loadedIdx.Lookup(hash1)
	if err != nil || len(positions1) != 1 || positions1[0] != pos1 {
		t.Errorf("Expected position %d for hash %x, got %v", pos1, hash1, positions1)
	}
	positions2, err := loadedIdx.Lookup(hash2)
	if err != nil || len(positions2) != 1 || positions2[0] != pos2 {
		t.Errorf("Expected postion %d for hash %x, got %v", pos2, hash2, positions2)
	}
}