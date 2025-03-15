package index

import (
	"sync"
	"time"
	"jamtext/internal/simhash"
)

// IndexShard represents a portion of the index
type IndexShard struct {
	SimHashToPos map[simhash.SimHash][]int64
	LSHBuckets   map[string]*LSHBucket // Added LSH support
	ShardID      int
	LastAccess   time.Time
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
	cachedShards  map[int]*IndexShard           // Cache for frequently accessed shards
	cacheSize     int                           // Maximum number of shards to keep in memory
	shardMap      map[simhash.SimHash]int       // Maps hashes to their shard IDs
}

// IndexStats contains statistics about the index
type IndexStats struct {
	TotalChunks    int64
	UniqueHashes   int64
	TotalPositions int64
	ShardCount     int
	MemoryUsage    int64
	CreationTime   time.Time
}
