# Jam-Text

A high-performance text indexer using SimHash fingerprints for text similarity search. Written in Go, it provides efficient indexing and searching of large text files through vector similarity with random hyperplanes.

## Current Status

✅ Implemented:
- Parallel chunk processing architecture
- SimHash core implementation
- Worker pool for concurrent processing
- CLI framework
- Basic project structure

🚧 In Progress:
- Index command implementation
- Lookup command implementation
- Index storage serialization
- Full CLI functionality

## Core Components

### Chunk Processing
- Default chunk size: 4KB
- Configurable overlap: 256 bytes
- Boundary-aware splitting
- Metadata support per chunk
- Parallel processing via worker pool

### SimHash Implementation
- 128-dimensional vector space
- 64-bit fingerprints
- Parallel hyperplane generation
- Box-Muller transform for normal distribution
- Normalized random hyperplanes
- Hamming distance similarity comparison

### Worker Pool
- Context-based graceful shutdown
- Buffered task channels
- Dynamic worker scaling
- Concurrent task processing

## Usage

```bash
# Index a file (in development)
jamtext -cmd index -i <input_file> -o <output_file>

# Lookup similar text (in development)
jamtext -cmd lookup -i <input_file> -o <output_file>
```

## Building

Requirements:
- Go 1.24.1 or higher
- Make (optional)

```bash
# Using make
make build

# Or manually
go build ./cmd/main.go
```

## Project Structure

```
.
├── cmd/
│   └── main.go          # Entry point
├── internal/
│   ├── cli/            # Command handling
│   ├── chunk/          # Text chunking
│   ├── index/          # Index management
│   └── simhash/        # SimHash implementation
├── go.mod
└── Makefile
```

## Technical Details

### Chunk Options
```go
ChunkOptions {
    ChunkSize:        4096,  // Default chunk size
    OverlapSize:      256,   // Overlap between chunks
    SplitOnBoundary:  true,  // Respect text boundaries
    BoundaryChars:    ".!?\n",
    MaxChunkSize:     6144,
    PreserveNewlines: true,
}
```

### SimHash Parameters
- Vector Dimensions: 128
- Number of Hyperplanes: 64
- Supported Vectorization Methods:
  - Frequency-based
  - N-gram based

### Performance Features
- Parallel hyperplane generation
- Concurrent chunk processing
- Buffered worker pools
- Context-based cancellation
- Efficient memory management

## Development

### Ignored Files
- `*.txt`: Text files
- `*.idx`: Index files
- `*.shard`: Shard data
- `*.dat`: Data files
- `.vscode/`: IDE settings

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

[License information to be added]
