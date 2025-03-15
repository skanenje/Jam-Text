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
***Syntax** <br>
```bash
./textindex -c index -i <input_file.txt> -s <chunk_size> -o <index_file.idx> [OPTIONAL= -log index.log]
```
**-c index** - Specifies that the command is for indexing the file <br>
**-i <input_file.txt>** - Path to the input text files <br>
**-s <chunk_size>** - Size of each chunk in bytes (default:4096 bytes) <br>
**-o <index_file.idx>** - Path to save the generated index file <br>
**-log <index.logs> - Path to the logs file to be able to test the hashes <br>
#### Example Usage
```
./textindex -c index -i large_text.txt -s 4096 -o index.idx -log index.log
```
### Looking up a chunk by Simhash
The **lookup** command retrieves the position of a chunk in a file based on its SimHash fingerprint
**Syntax** <br>
```bash
./textindex -c lookup -i <index.file.idx> -h <simhash_value>
```
**-c lookup** - Specifies that the command is for looking up a chunk <br>
**-i <index_file.idx>** - Path to the previously generated index file <br>
**-s <simhash_value>** - Simhash value of the chunk to search for <br>
#### Example Usage
```
./textindex -c lookup -i index_file.idx -h 3eff1b2c98a6
```
### Compare two documents similarity
The command **compare** compares two text documents and returns how much similar the documents are in percentage.
**Syntax** <br>
```bash
# compare two text documents to find if they are similar
./textindex -c compare -i doc1.txt -i2 doc2.txt -o report.txt
```
**-c compare** - Specifies that the command is for comparing two text documents <br>
**-i <input_file>** - Path to the first input document <br>
**-i2 <input_file2>** - Path to the second input document <br>
**-0 <report_file>** - Path to the report document that show the output. (OPTIONAl) <br>
#### Example Usage
```
./textindex -c compare -i file1.txt -i2 file2.txt -o report.txt
```
### Duplicate Detection
**1. index the desired corpus of data** <br>
**Syntax** <br>
```bash
./textindex -c index -i testdata.txt -o testdata.dat -s 1024 -overlap 256
```
**-c index** - Specifies that the command is for indexing the file <br>
**-i testdata.txt** - specifies path to text file <br>
**-o testdata.idx** -  Path to the previously generated index file <br>
**-s 1024** - Specifies the chunk size <br>
**overlap 256** - Specifies the amount of words to add at the beginning of the chunk and at the end <br>

**2.  hash the particular document you want** <br>
**Syntax** <br>
```
HASH=$(./textindex -c hash -i testPlagurism.txt)
```
**3. Use the hash for the lookup** <br>
**Syntax** <br>
```bash
./textindex -c fuzzy -i testdata.idx -h $HASH -threshold 5
```
**-c fuzzy** - This command specifies that you are looking for exact matches <br>
**-i testdata.idx** - path to the hash file generated after hashing a document <br>
**-h $HASH** - specifies the the hash <br>
**-threshold -5** - specifies the limiting values <br> 

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
- [Testing](docs/testing.md) - Testing practices and patterns.

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
We welcome contributions! Please see [CONTRIBUTING.md](docs/CONTRIBUTING.md).

## Testing
Run the test suite using the following commands:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/cli/...
go test ./internal/simhash/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributors
- [Elijah Gathanga](https://github.com/GathangaElijah)
- [Kevin Wasonga](https://github.com/kevwasonga)
- [Swabri Kanenje](https://github.com/skanenje)
- [Jerome Otieno](https://github.com/Jerome-afk) 
- [Godwin Ouma](https://github.com/oumaoumag) 

## License
MIT License - see [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE) for details.
