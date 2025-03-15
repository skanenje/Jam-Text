# Jam-Text

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24.1-blue.svg)](https://golang.org/dl/)

A high-performance text indexer using SimHash fingerprints for text similarity search and content finding.

## Features
- SimHash-based text fingerprinting
- Efficient chunk processing with configurable sizes
- Parallel processing with worker pools
- Locality-Sensitive Hashing (LSH) support
- Multiple vectorization strategies
- Content finding and search capabilities
- Fuzzy matching with configurable thresholds
- Sharded index storage for large datasets
- Memory-efficient operations through disk-based sharding

## Quick Start
```bash
# Build a executable of the program
make

# Index a document
./textindex -c index -i testdata.txt -o testdata.dat -s 1024

# Generate document hash
HASH=$(textindex -c hash -i testPlagiarism.txt)

# Search with fuzzy matching
./textindex -c fuzzy -i testdata.dat -h $HASH -threshold 5

# Find known content
KNOWN_HASH=$(textindex -c hash -i known_content.txt)
./textindex -c lookup -i database.idx -h $KNOWN_HASH
```

## Common Use Cases
- **Content Finding**: Search for known text across large document collections
- **Plagiarism Detection**: Compare documents for similarity
- **Content Moderation**: Screen content against moderation rules
- **Database Search**: Build and search text databases efficiently
- **Document Comparison**: Compare multiple documents for similarities

## Core Components

### Index System
The Index package provides the core functionality for text indexing and search:
- Sharded storage supporting 100,000+ entries per shard
- Automatic shard rotation and management
- LSH-based similarity search
- Thread-safe operations
- Memory-efficient disk-based storage

For detailed index implementation, see [Index Documentation](docs/Index_Readme.md).

## Package Structure
- `cmd/` - Command-line interface
- `internal/` - Internal implementation
  - `chunk/` - Text segmentation
  - `simhash/` - Fingerprint generation
  - `cli/` - CLI implementation
  - `index/` - Core indexing system
- `docs/` - Package documentation

## Documentation
See package documentation for detailed information:
- [Index System](docs/Index_Readme.md) - Core indexing and search functionality
- [CLI Documentation](docs/cli.md) - Command-line interface
- [Chunk Package](docs/chunking.md) - Text segmentation
- [SimHash Package](docs/simhash.md) - Fingerprint generation

## Search Examples
```bash
# Create content database with efficient sharding
./textindex -c index -i master_content.txt -o reference.idx -s 4096

# Search for known content
HASH=$(textindex -c hash -i known_text.txt)
./textindex -c fuzzy -i reference.idx -h $HASH -threshold 3

# Batch search multiple indexes
for idx in indexes/*.idx; do
./textindex -c fuzzy -i "$idx" -h $HASH -threshold 3
done
```

## Performance Tips
- Configure shard sizes based on available memory (default: 100,000 entries)
- Use LSH bands appropriately for your dataset size
- Adjust chunk sizes:
  - Smaller (2048) for precise matching
  - Larger (8192+) for faster processing
- Monitor shard rotation frequency for optimal performance

## Requirements
- Go 1.24.1 or higher
- External dependencies:
  - `pdftotext` (poppler-utils) for PDF support
  - `pandoc` for DOCX support

## Contributing
We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License
MIT License - see [LICENSE](LICENSE) for details.
