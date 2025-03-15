# Jam-Text Documentation

## Core Documentation
- [CLI Reference](cli.md) - Command-line interface documentation
  - Commands: index, lookup, fuzzy, hash, compare, moderate
  - Usage examples and configuration options
  - Advanced parameters for chunk size, overlap, and thresholds

## Technical Documentation
- [SimHash Implementation](simhash.md) - Text fingerprinting and similarity detection
  - 64-bit fingerprint representation
  - Frequency-based and n-gram vectorization
  - LSH support for fast similarity search
  - Thread-safe operations
- [Chunking System](chunking.md) - Text segmentation and processing
  - Text segmentation with configurable chunk sizes
  - Context-aware processing with overlap support
  - Multi-format support (TXT, PDF, DOCX)
  - Parallel processing with worker pools
- [Testing Guide](testing.md) - Comprehensive testing documentation
  - Test structure and organization
  - Testing techniques and patterns
  - Package-specific test documentation
  - Best practices and examples
  - CI/CD testing guidelines

## Contributing
- [Contributing Guide](../CONTRIBUTING.md) - Guidelines for contributors
