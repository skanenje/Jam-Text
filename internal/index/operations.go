package index

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"jamtext/internal/simhash"
)

const (
	MaxShardSize    = 100000
	ShardTimeoutMin = 30
)

// New creates a new Index
func New(sourceFile string, chunkSize int, hyperplanes [][]float64, indexDir string) *Index {
	if indexDir == "" {
		indexDir = filepath.Join(os.TempDir(), "textindex")
	}

	os.MkdirAll(indexDir, 0o755)

	return &Index{
		SourceFile:    sourceFile,
		ChunkSize:     chunkSize,
		Hyperplanes:   hyperplanes,
		CreationTime:  time.Now(),
		LSHTable:      simhash.NewPermutationTable(simhash.NumHyperplanes, 4),
		IndexDir:      indexDir,
		ShardFilename: filepath.Base(sourceFile) + ".shard",
		Shards: []*IndexShard{{
			SimHashToPos: make(map[simhash.SimHash][]int64),
			ShardID:      0,
			LastAccess:   time.Now(),
		}},
	}
}

// Add adds a SimHash and position to the index with LSH support
func (idx *Index) Add(hash simhash.SimHash, pos int64) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Add to regular index
	shard := idx.Shards[idx.ActiveShard]
	shard.SimHashToPos[hash] = append(shard.SimHashToPos[hash], pos)

	// Add to LSH buckets
	signatures := idx.LSHTable.GetBandSignatures(hash)
	for i, sig := range signatures {
		bucketKey := fmt.Sprintf("%d:%d", i, sig)
		if shard.LSHBuckets == nil {
			shard.LSHBuckets = make(map[string]*LSHBucket)
		}
		if shard.LSHBuckets[bucketKey] == nil {
			shard.LSHBuckets[bucketKey] = &LSHBucket{
				hashes: make(map[simhash.SimHash]struct{}),
			}
		}
		shard.LSHBuckets[bucketKey].hashes[hash] = struct{}{}
	}

	if len(shard.SimHashToPos) >= MaxShardSize {
		if err := idx.rotateShard(); err != nil {
			return fmt.Errorf("failed to rotate shard: %w", err)
		}
	}

	return nil
}

// saveShard persists a shard to disk
func (idx *Index) saveShard(shard *IndexShard) error {
	filename := filepath.Join(idx.IndexDir, idx.ShardFilename+"."+string(rune('0'+shard.ShardID)))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(shard.SimHashToPos)
}

// loadShard loads a shard from disk
func (idx *Index) loadShard(shardID int) (*IndexShard, error) {
	filename := filepath.Join(idx.IndexDir, idx.ShardFilename+"."+string(rune('0'+shardID)))
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var simHashToPos map[simhash.SimHash][]int64
	if err := gob.NewDecoder(file).Decode(&simHashToPos); err != nil {
		return nil, err
	}

	return &IndexShard{
		SimHashToPos: simHashToPos,
		ShardID:      shardID,
	}, nil
}

// rotateShard persists the current shard and creates a new one
func (idx *Index) rotateShard() error {
	if err := idx.saveShard(idx.Shards[idx.ActiveShard]); err != nil {
		return err
	}

	idx.ActiveShard++
	idx.Shards = append(idx.Shards, &IndexShard{
		SimHashToPos: make(map[simhash.SimHash][]int64),
		ShardID:      idx.ActiveShard,
		LastAccess:   time.Now(),
	})

	return nil
}

// Lookup finds positions for a SimHash
func (idx *Index) Lookup(hash simhash.SimHash) ([]int64, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var positions []int64

	// Check active shard first
	if pos, ok := idx.Shards[idx.ActiveShard].SimHashToPos[hash]; ok {
		positions = append(positions, pos...)
	}

	// Check other shards
	for i := 0; i < len(idx.Shards); i++ {
		if i == idx.ActiveShard {
			continue
		}

		shard := idx.Shards[i]
		if shard == nil {
			// Load shard from disk
			loaded, err := idx.loadShard(i)
			if err != nil {
				return nil, fmt.Errorf("failed to load shard %d: %w", i, err)
			}
			idx.Shards[i] = loaded
			shard = loaded
		}

		if pos, ok := shard.SimHashToPos[hash]; ok {
			positions = append(positions, pos...)
		}
	}

	return positions, nil
}

// Stats returns statistics about the index
func (idx *Index) Stats() map[string]interface{} {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	totalEntries := 0
	totalPositions := 0

	for _, shard := range idx.Shards {
		if shard != nil {
			totalEntries += len(shard.SimHashToPos)
			for _, positions := range shard.SimHashToPos {
				totalPositions += len(positions)
			}
		}
	}

	return map[string]interface{}{
		"source_file":     idx.SourceFile,
		"chunk_size":      idx.ChunkSize,
		"created":         idx.CreationTime,
		"shards":          len(idx.Shards),
		"unique_hashes":   totalEntries,
		"total_positions": totalPositions,
	}
}

// Close performs cleanup operations
func (idx *Index) Close() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Save active shard if needed
	if len(idx.Shards[idx.ActiveShard].SimHashToPos) > 0 {
		if err := idx.saveShard(idx.Shards[idx.ActiveShard]); err != nil {
			return err
		}
	}

	// Clear in-memory data to free up resources
	for i := range idx.Shards {
		idx.Shards[i] = nil
	}

	return nil
}

// FuzzyLookup finds positions for similar SimHashes using LSH
func (idx *Index) FuzzyLookup(hash simhash.SimHash, threshold int) (map[simhash.SimHash][]int64, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	candidates := make(map[simhash.SimHash]struct{})
	signatures := idx.LSHTable.GetBandSignatures(hash)

	// Collect candidates from LSH buckets
	for i, sig := range signatures {
		bucketKey := fmt.Sprintf("%d:%d", i, sig)
		for _, shard := range idx.Shards {
			if shard == nil || shard.LSHBuckets == nil {
				continue
			}
			if bucket := shard.LSHBuckets[bucketKey]; bucket != nil {
				for candidateHash := range bucket.hashes {
					candidates[candidateHash] = struct{}{}
				}
			}
		}
	}

	// Verify candidates with Hamming distance
	results := make(map[simhash.SimHash][]int64)
	found := false

	for candidateHash := range candidates {
		if candidateHash.IsSimilar(hash, threshold) {
			for _, shard := range idx.Shards {
				if shard == nil {
					continue
				}
				if positions, ok := shard.SimHashToPos[candidateHash]; ok {
					results[candidateHash] = append(results[candidateHash], positions...)
					found = true
				}
			}
		}
	}

	return results, found
}

// LSHBucket represents a collection of similar hashes
type LSHBucket struct {
	hashes map[simhash.SimHash]struct{}
}
