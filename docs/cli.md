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
The CLI supports various flags for configuration:

##### Basic Flags
- `-c`: Command to execute (required)
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

#### Example
```bash
jamtext -c index -i book.txt -o book.idx -s 4096 -overlap 256
```

### 2. Lookup Command
Searches indexed files for similar text chunks.

#### Usage
```bash
jamtext -c lookup -i <index_file> -h <hash> [options]
```

#### Operation Flow
1. **Index Loading**
   - Opens and validates index file
   - Prepares search structures

2. **Search Process**
   - Locates matching positions
   - Retrieves original text chunks
   - Adds contextual information

3. **Results Display**
   - Shows top matches (maximum 3)
   - Includes position and text preview

#### Example
```bash
jamtext -c lookup -i book.idx -h f7a3d921 -context-before 100 -context-after 100
```

## Integration with Other Packages

### 1. Chunk Package (`internal/chunk/`)
- Text segmentation management
- Chunk processing coordination
- File reading operations

### 2. Index Package (`internal/index/`)
- Fingerprint storage
- Index persistence
- Search operations

### 3. SimHash Package (`internal/simhash/`)
- Text fingerprinting
- Similarity calculations
- Hyperplane management

## Error Handling

The CLI package implements comprehensive error handling:

1. **Input Validation**
   - Required parameter checking
   - File path validation
   - Command syntax verification

2. **Runtime Error Management**
   - File operation errors
   - Processing failures
   - Resource allocation issues

3. **User Feedback**
   - Clear error messages
   - Processing statistics
   - Operation status updates

## Logging

### Configuration
- Default: Logs to stderr
- Optional file logging with `-log` flag
- Verbose mode with `-v` flag

### Log Content
- Operation progress
- Error information
- Performance metrics
- Debug information (in verbose mode)

## Performance Considerations

1. **Resource Management**
   - Efficient file handling
   - Memory-conscious processing
   - Proper resource cleanup

2. **Processing Optimization**
   - Parallel processing support
   - Configurable chunk sizes
   - Optimized search operations

## Best Practices

1. **Index Command**
   - Use appropriate chunk sizes for your content
   - Enable boundary splitting for text files
   - Configure overlap based on content structure

2. **Lookup Command**
   - Provide sufficient context for meaningful results
   - Use appropriate similarity thresholds
   - Consider index size for performance

## Examples

### Basic Indexing
```bash
# Index with default settings
jamtext -c index -i document.txt -o document.idx

# Index with custom chunk size
jamtext -c index -i document.txt -o document.idx -s 8192
```

### Advanced Indexing
```bash
# Index with custom settings
jamtext -c index -i document.txt -o document.idx \
    -s 4096 \
    -overlap 512 \
    -boundary=true \
    -boundary-chars ".!?\n" \
    -max-size 8192 \
    -preserve-nl=true \
    -index-dir ./indexes
```

### Lookup Operations
```bash
# Basic lookup
jamtext -c lookup -i document.idx -h <hash>

# Lookup with extended context
jamtext -c lookup -i document.idx -h <hash> \
    -context-before 200 \
    -context-after 200
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