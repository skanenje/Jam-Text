# CLI Package Documentation

## Overview

The Command Line Interface (CLI) package serves as the central control system for JamText, managing user interactions and coordinating program operations. This document details the CLI package's architecture, functionality, and integration with other system components.

## Core Components

### 1. Entry Point (`cmd/main.go`)
The program's entry point forwards command-line arguments to the CLI package:
```go
func main() {
    if err := cli.Run(os.Args); err != nil {
        panic(err)
    }
}
```

### 2. Command Processing (`internal/cli/cli.go`)

#### Command Line Options

##### Basic Flags
- `-c`: Command to execute (required)
  - `index`: Create index from text file
  - `lookup`: Exact lookup by SimHash
  - `fuzzy`: Fuzzy lookup with threshold
  - `hash`: Calculate SimHash for a file
  - `stats`: Show index statistics
  - `compare`: Compare two text files
- `-i`: Input file path
- `-o`: Output file path
- `-v`: Enable verbose output
- `-log`: Log file path (defaults to stderr)
- `-s`: Chunk size in bytes (default: 4096)
- `-h`: SimHash value for lookup

##### Advanced Flags
- `-overlap`: Overlap size in bytes (default: 256)
- `-boundary`: Enable text boundary splitting (default: true)
- `-boundary-chars`: Characters for text boundaries (default: ".!?\n")
- `-max-size`: Maximum chunk size in bytes (default: 6144)
- `-preserve-nl`: Preserve newlines in chunks (default: true)
- `-index-dir`: Directory for index shard storage
- `-context-before`: Context bytes before chunk (default: 100)
- `-context-after`: Context bytes after chunk (default: 100)

## Main Commands

### 1. Index Command
Processes and indexes text files for similarity searching.

#### Usage
```bash
jamtext -c index -i <input_file> -o <output_file> [options]
```

#### Operation Flow
1. **Preparation**
   - Validates input/output paths
   - Generates SimHash hyperplanes
   - Configures chunking options

2. **Processing**
   - Splits input file into chunks
   - Generates fingerprints
   - Creates searchable index

3. **Finalization**
   - Saves index to specified output
   - Displays processing statistics

### 2. Compare Command
Compares two text files for similarity.

#### Usage
```bash
jamtext -c compare -i <file1> -i2 <file2> [-o <report_file>]
```

#### Operation Flow
1. **Document Loading**
   - Reads both input files
   - Initializes similarity detector

2. **Comparison Process**
   - Calculates SimHash fingerprints
   - Performs similarity analysis
   - Generates detailed comparison report

3. **Output**
   - Displays similarity metrics
   - Optionally saves detailed report

### 3. Hash Command
Generates SimHash for a single document.

#### Usage
```bash
jamtext -c hash -i <input_file>
```

#### Operation Flow
1. **Processing**
   - Reads input file
   - Generates hyperplanes
   - Calculates SimHash

2. **Output**
   - Displays 64-bit hash value in hexadecimal format (%x)
   - Returns error if hash parsing fails

### 4. Fuzzy Command
Performs fuzzy lookup using SimHash with configurable threshold.

#### Usage
```bash
jamtext -c fuzzy -i <index_file> -h <hash> -threshold <value>
```

#### Operation Flow
1. **Hash Parsing**
   - Parses input hash from hexadecimal string
   - Validates hash format
2. **Lookup Process**
   - Performs LSH-enhanced fuzzy lookup
   - Applies similarity threshold
3. **Output**
   - Shows matching chunks with context
   - Displays similarity metrics

## Integration with Other Packages

### Package Dependencies
- `chunk`: Text segmentation
- `simhash`: Fingerprint generation
- `index`: Storage management

### Error Handling
- Comprehensive input validation
- Detailed error reporting
- Resource cleanup on failure

## Best Practices

### Command Selection
- Use `index` for corpus preparation
- `compare` for direct file comparison
- `hash` for single document fingerprinting
- `fuzzy` for similarity search

### Performance Optimization
- Choose appropriate chunk sizes
- Enable boundary splitting for text
- Configure overlap based on content

## Examples

### Basic Operations
```bash
# Index a document
jamtext -c index -i book.txt -o book.idx

# Generate document hash
jamtext -c hash -i document.txt

# Compare two files
jamtext -c compare -i doc1.txt -i2 doc2.txt -o report.txt
```

### Advanced Usage
```bash
# Index with custom settings
jamtext -c index -i book.txt -o book.idx \
    -s 4096 -overlap 512 -boundary=true

# Fuzzy search with context
jamtext -c fuzzy -i book.idx -h <hash> \
    -context-before 200 -context-after 200
```

## Future Enhancements

1. **Planned Features**
   - Additional vectorization methods
   - Enhanced similarity metrics
   - Improved chunk boundary detection

2. **Performance Optimizations**
   - Enhanced parallel processing
   - Improved memory management
   - Optimized search algorithms

## Support and Maintenance

For issues, suggestions, or contributions:
1. Check existing documentation
2. Review error messages
3. Consult the project repository
4. Submit detailed bug reports