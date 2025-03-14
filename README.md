# Jam-Text

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24.1-blue.svg)](https://golang.org/dl/)

[Jam-Text](https://github.com/yourusername/jam-text) is a high-performance text indexer using SimHash fingerprints for text similarity search. Written in Go, it provides efficient indexing and searching of large text files through vector similarity with random hyperplanes.

## Installing

For the latest stable version:

```bash
go install github.com/yourusername/jam-text@latest
```

Or build from source:

```bash
make build
```

## Quick Start

Index a document and search for similar text:

```bash
# Index the desired corpus of data
jamtext -c index -i testdata.txt -o testdata.dat -s 1024 -overlap 256

# Hash a particular document
HASH=$(jamtext -c hash -i testPlagiarism.txt)

# Use the hash for lookup with fuzzy search
jamtext -c fuzzy -i testdata.dat -h $HASH -threshold 5
```

## Contribute

There are many ways to contribute to Jam-Text:
* [Submit bugs](https://github.com/yourusername/jam-text/issues) and help verify fixes
* Review [source code changes](https://github.com/yourusername/jam-text/pulls)
* [Contribute bug fixes](CONTRIBUTING.md)

## Documentation

* [Getting Started](docs/guides/getting-started.md)
* [CLI Documentation](docs/cli.md)
* [Architecture Guide](docs/architecture/design.md)
* [API Reference](docs/api/README.md)

## Core Features

* SimHash-based text fingerprinting
* Efficient chunk processing with configurable sizes
* Parallel processing with worker pools
* Locality-Sensitive Hashing (LSH) support
* Two vectorization strategies:
  * Frequency-based with MD5 dimension mapping
  * N-gram based with normalized vectors

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
│   └── simhash/        # SimHash implementation
├── go.mod
└── Makefile
```

## Roadmap

For details on planned features and future direction, please refer to our [roadmap](docs/roadmap.md).
