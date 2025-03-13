package index

import (
	"sync"
	"time"

	"jamtext/internal/simhash"
)

// Maximum number of entries per shard
const (
	MaxShardSize = 100000
)

// IndexShard represents a portion of the index
type IndexShard struct {
	SimHashToPos map[simhash.SimHash][]int64
	ShardID      int
}

// Index stores SimHash mappings with sharding support
type Index struct {
	SourceFile    string
	ChunkSize     int
	Shards        []*IndexShard
	ActiveShard   int
	Hyperplanes   [][]float64
	CreationTime  time.Time
	LSHTable      *simhash.PermutationTable
	IndexDir      string
	mu            sync.RWMutex
	ShardFilename string
}
