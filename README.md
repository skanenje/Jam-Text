# Jam-Text

A high-performance text indexer implementation using SimHash fingerprints based on vector similarity with random hyperplanes. Written in Go, it provides an efficient way to index and search large text files.

## How It Works

The text indexer uses a combination of techniques for efficient text similarity search:

1. **Chunking**: Text files are split into fixed-size chunks (default 4KB) using a parallel worker pool for performance
2. **Vectorization**: Each chunk is converted to a normalized word frequency vector using either:
   - Simple frequency-based vectorization
   - N-gram based vectorization
3. **SimHash Generation**: 
   - Uses 64-bit SimHash fingerprints
   - Implements random hyperplane generation for consistent hashing
   - Employs parallel processing for hyperplane generation
4. **Similarity Detection**:
   - Uses Hamming distance for comparing SimHash fingerprints
   - Supports configurable similarity thresholds
   - Includes locality-sensitive hashing for faster similarity search

## Features

- Multi-threaded processing using worker pools for chunk processing
- Vector-based SimHash computation using 128-dimensional vectors
- 64 random hyperplanes for fingerprint generation
- Efficient memory use with normalized word-frequency vectors
- Clean command-line interface
- Serialized index storage using Gob encoding

## Building the Application

Requirements:
- Go 1.24.1 or higher
- Make (optional)

To build the application:

```bash
make build
```

Or manually:

```bash
go build ./cmd/main.go
```

## Usage

The application supports two main commands:

### Indexing
```bash
jamtext -cmd index -i <input_file> -o <output_file>
```

### Lookup
```bash
jamtext -cmd lookup -i <input_file> -o <output_file>
```

## Project Structure

```
.
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── cli/            # Command line interface handling
│   ├── chunk/          # Parallel chunk processing implementation
│   ├── index/          # Index management and storage
│   └── simhash/        # SimHash implementation with LSH support
├── go.mod              # Go module definition
└── Makefile            # Build automation
```

## Technical Details

### Chunk Processing
- Implements a worker pool for parallel chunk processing
- Uses context for graceful shutdown
- Includes metadata support for each chunk
- Buffered channels for optimal performance

### SimHash Implementation
- 128-dimensional vector space
- 64-bit fingerprint generation
- Parallel hyperplane generation using Box-Muller transform
- Normalized random hyperplanes for consistent hashing
- Supports both frequency and n-gram based vectorization

### Index Management
- In-memory index mapping fingerprints to file positions
- Serialized storage support
- Efficient similarity search capabilities

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[License information to be added]
