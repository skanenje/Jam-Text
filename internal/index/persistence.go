package index

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"jamtext/internal/simhash"
)

// Save writes the index metadata to a file
func Save(idx *Index, outputFile string) error {
	// First save any active shard
	if err := idx.saveShard(idx.Shards[idx.ActiveShard]); err != nil {
		return fmt.Errorf("failed to save active shard: %w", err)
	}

	// Create metadata structure
	meta := struct {
		SourceFile    string
		ChunkSize     int
		ShardCount    int
		Hyperplanes   [][]float64
		CreationTime  time.Time
		IndexDir      string
		ShardFilename string
	}{
		SourceFile:    idx.SourceFile,
		ChunkSize:     idx.ChunkSize,
		ShardCount:    len(idx.Shards),
		Hyperplanes:   idx.Hyperplanes,
		CreationTime:  idx.CreationTime,
		IndexDir:      idx.IndexDir,
		ShardFilename: idx.ShardFilename,
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(meta); err != nil {
		return fmt.Errorf("failed to encode index metadata: %w", err)
	}

	return nil
}

// Load reads an index from a file
func Load(indexFile string) (*Index, error) {
	file, err := os.Open(indexFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}
	defer file.Close()

	var meta struct {
		SourceFile    string
		ChunkSize     int
		ShardCount    int
		Hyperplanes   [][]float64
		CreationTime  time.Time
		IndexDir      string
		ShardFilename string
	}

	if err := gob.NewDecoder(file).Decode(&meta); err != nil {
		return nil, fmt.Errorf("failed to decode index metadata: %w", err)
	}

	idx := &Index{
		SourceFile:    meta.SourceFile,
		ChunkSize:     meta.ChunkSize,
		Hyperplanes:   meta.Hyperplanes,
		CreationTime:  meta.CreationTime,
		LSHTable:      simhash.NewPermutationTable(simhash.NumHyperplanes, 4),
		IndexDir:      meta.IndexDir,
		ShardFilename: meta.ShardFilename,
		Shards:        make([]*IndexShard, meta.ShardCount),
	}

	// Load first shard
	if meta.ShardCount > 0 {
		firstShard, err := idx.loadShard(0)
		if err != nil {
			return nil, fmt.Errorf("failed to load first shard: %w", err)
		}
		idx.Shards[0] = firstShard
	}

	return idx, nil
}
