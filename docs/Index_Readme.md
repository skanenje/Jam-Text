
# Index Package Documentation

## Overview
The Index package provides core functionality for JamText's text indexing and search capabilities. It implements:
- Sharded index storage
- LSH (Locality-Sensitive Hashing) support
- Efficient SimHash lookup
- Persistent storage management

## Core Components

### Index Structure
```go
type Index struct {
    SourceFile    string
    ChunkSize     int
    Shards        []*IndexShard
    ActiveShard   int
    Hyperplanes   [][]float64
    CreationTime  time.Time
    LSHTable      *simhash.PermutationTable
    IndexDir      string
}
```

### Key Features
- Automatic shard rotation (MaxShardSize: 100,000 entries)
- Memory-efficient operation through disk-based sharding
- Thread-safe operations
- LSH-based similarity search

## Usage Examples

### Creating an Index
```go
// Initialize new index
idx := index.New(sourceFile, chunkSize, hyperplanes, indexDir)

// Add content
idx.Add(hash, position)

// Save index
index.Save(idx, outputPath)
```

### Search Operations
```go
// Exact lookup
positions, err := idx.Lookup(hash)

// Fuzzy search
matches, found := idx.FuzzyLookup(hash, threshold)
```

### Shard Management
```go
// Save active shard
idx.saveShard(activeShard)

// Load specific shard
shard, err := idx.loadShard(shardID)
```

## Performance Considerations
- Default shard size: 100,000 entries
- Shard timeout: 30 minutes
- LSH configuration affects search speed vs accuracy
- Use appropriate chunk sizes for your use case

## Best Practices
1. Configure shard sizes based on available memory
2. Implement regular shard cleanup
3. Use appropriate LSH bands for dataset size
4. Monitor shard rotation frequency

## Integration with Other Packages
- Works with `simhash` package for fingerprint generation
- Integrates with `chunk` package for text processing
- Supports CLI operations through `cli` package

## Example Workflows

### Building Search Index
```go
idx := index.New(sourceFile, 4096, hyperplanes, "/tmp/index")
for _, chunk := range chunks {
    hash := simhash.Calculate(chunk)
    idx.Add(hash, chunk.Position)
}
index.Save(idx, "output.idx")
```

### Implementing Search
```go
idx, err := index.Load("database.idx")
if err != nil {
    return err
}
matches, found := idx.FuzzyLookup(targetHash, 3)
```

For detailed implementation examples, see the test files and CLI documentation.