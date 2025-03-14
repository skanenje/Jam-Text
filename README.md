# Jam-Text

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24.1-blue.svg)](https://golang.org/dl/)

A high-performance text indexer using SimHash fingerprints for text similarity search. Written in Go, it provides efficient indexing and searching of large text files through vector similarity with random hyperplanes.

## Key Features

- SimHash-based text fingerprinting
- Efficient chunk processing with configurable sizes
- Parallel processing with worker pools
- Locality-Sensitive Hashing (LSH) support
- Multiple vectorization strategies:
  - Frequency-based with MD5 dimension mapping
  - N-gram based with normalized vectors

## Quick Start

### Installation

```bash
# Latest stable version
go install github.com/yourusername/jam-text@latest

# Or build from source
make build
```

### Basic Usage

```bash
# Index a document
jamtext -c index -i testdata.txt -o testdata.dat -s 1024 -overlap 256

# Generate document hash
HASH=$(jamtext -c hash -i testPlagiarism.txt)

# Search with fuzzy matching
jamtext -c fuzzy -i testdata.dat -h $HASH -threshold 5
```

## Documentation

- [Getting Started Guide](docs/guides/getting-started.md)
- [CLI Documentation](docs/cli.md)
- [Architecture Overview](docs/architecture/design.md)
- [API Reference](docs/api/README.md)

## Requirements

- Go 1.24.1 or higher
- Make (optional)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.
