# Jam-Text

A high-performance text indexer using SimHash fingerprints for text similarity search. Written in Go, it provides efficient indexing and searching of large text files through vector similarity with random hyperplanes.
## Core Components

### SimHash Implementation
- 128-dimensional vector space
- 64-bit fingerprints
- Parallel hyperplane generation using Box-Muller transform
- Normalized random hyperplanes
- Locality-Sensitive Hashing (LSH) support
- Configurable band signatures for fast similarity search
- Two vectorization strategies:
  - Frequency-based: Uses word frequencies with MD5 dimension mapping
  - N-gram based: Uses character n-grams with normalized vectors

### Chunk Processing
- Default chunk size: 4KB
- Configurable overlap: 256 bytes
- Boundary-aware splitting
- Metadata support per chunk
- Parallel processing via worker pool

### Worker Pool
- Context-based graceful shutdown
- Buffered task channels
- Dynamic worker scaling
- Concurrent task processing

## Use Cases

## Plugerism Detection
```bash
# index the desired corpus of data
./jamtext -c index -i testdata.txt -o testdata.dat -s 1024 -overlap 256

# hash the particular document you want
HASH=$(./jamtext -c hash -i testPlagurism.txt)

# use the hash for lookup with fuzzy search
./jamtext -c fuzzy -i testdata.dat -h $HASH -threshold 5
```

## Similarity Search
```bash
# compare two text documents to find if they are similar
./jamtext -c compare -i doc1.txt -i2 doc2.txt -o report.txt

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
│   └── simhash/        # SimHash implementation with LSH
├── go.mod
└── Makefile
```

## Technical Details

### SimHash Parameters
```go
const (
    VectorDimensions = 128
    NumHyperplanes   = 64
)
```

### Vectorization Options
- Frequency-based vectorization:
  - Word-level tokenization
  - MD5-based dimension mapping
  - Vector normalization
  - Thread-safe implementation
- N-gram vectorization:
  - Configurable n-gram size
  - Normalized vector output
  - Fallback for short texts
  - Efficient n-gram generation

### LSH Configuration
- Configurable band size
- Random permutation generation
- Band signature computation
- Optimized for similarity search

### Performance Features
- Parallel hyperplane generation
- Concurrent chunk processing
- Buffered worker pools
- Context-based cancellation
- Efficient memory management
- Thread-safe random number generation
- Optimized vector operations

## Documentation

Comprehensive documentation is available in the [docs](docs/) directory:

- Package documentation with examples and best practices
- Architecture and design documents
- Performance tuning guides
- API reference

For package-level documentation, use `go doc`:

```bash
go doc jamtext/internal/simhash
go doc jamtext/internal/chunk
go doc jamtext/internal/index
```

## Development

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

[License information to be added]
