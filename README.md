# Jam-Text

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24.1-blue.svg)](https://golang.org/dl/)

A high-performance text indexer using SimHash fingerprints for text similarity search.

## Features
- SimHash-based text fingerprinting
- Efficient chunk processing with configurable sizes
- Parallel processing with worker pools
- Locality-Sensitive Hashing (LSH) support
- Multiple vectorization strategies

## Quick Start
```bash
# Install
go install github.com/yourusername/jam-text@latest

# Index a document
jamtext -c index -i testdata.txt -o testdata.dat -s 1024

# Generate document hash
HASH=$(jamtext -c hash -i testPlagiarism.txt)

# Search with fuzzy matching
jamtext -c fuzzy -i testdata.dat -h $HASH -threshold 5
```

## Package Structure
- `cmd/` - Command-line interface
- `internal/` - Internal implementation
  - `chunk/` - Text segmentation
  - `simhash/` - Fingerprint generation
  - `cli/` - CLI implementation
- `docs/` - Package documentation

## Documentation
See package READMEs for detailed documentation:
- [CLI Documentation](internal/cli/README.md)
- [Chunk Package](internal/chunk/README.md)
- [SimHash Package](internal/simhash/README.md)

## Requirements
- Go 1.24.1 or higher
- External dependencies:
  - `pdftotext` (poppler-utils) for PDF support
  - `pandoc` for DOCX support

## Contributing
We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License
MIT License - see [LICENSE](LICENSE) for details.
