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

### Step 1 Extraction
```
Extract the Jam-Text-main.zip
```
### Step 2 Opening the folder
```
Open Jam-Text-main in terminal
Rightclick and open in terminal
```
### Step 3 Build Executable
```bash
#run make to build the executable file of the program
make
```
### Step 4 Test the program
### Indexing a text file
The **index** command processes a text file and creates an in memory index of simhash values
***Syntax**
```bash
./textindex -c index -i <input_file.txt> -s <chunk_size> -o <index_file.idx> [OPTIONAL= -log index.log]
```
**-c index** - Specifies that the command is for indexing the file
**-i <input_file.txt>** - Path to the input text files
**-s <chunk_size>** - Size of each chunk in bytes (default:4096 bytes)
**-o <index_file.idx>** - Path to save the generated index file
**-log <index.logs> - Path to the logs file to be able to test the hashes 
#### Example Usage
```
./textindex -c index -i large_text.txt -s 4096 -o index.idx -log index.log
```
### Looking up a chunk by Simhash
The **lookup** command retrieves the position of a chunk in a file based on its SimHash fingerprint
**Syntax**
```bash
./textindex -c lookup -i <index.file.idx> -h <simhash_value>
```
**-c lookup** - Specifies that the command is for looking up a chunk
**-i <index_file.idx>** - Path to the previously generated index file
**-s <simhash_value>** - Simhash value of the chunk to search for
#### Example Usage
```
./textindex -c lookup -i index_file.idx -h 3eff1b2c98a6
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
