# Jam-Text

A high-performance text indexer implementation using SimHash fingerprints based on vector similarity with random hyperplanes. Written in Go, it provides an efficient way to index and search large text files.

## The Text Indexer Works By

1. Splitting text files into fixed-size chunks (default 4KB)
2. Converting each chunk to a normalized word frequency vector
3. Using random hyperplanes to generate SimHash fingerprints
4. Building an in-memory index that maps fingerprints to file positions
5. Allowing quick lookups based on hash values

## Features

- Multi-threaded processing for faster indexing
- Vector-based SimHash computation using random hyperplanes
- Efficient memory use with normalized word-frequency vectors
- Clean command-line interface matching the specification
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
│   ├── chunk/          # Text chunking implementation
│   ├── index/          # Index management
│   └── simhash/        # SimHash implementation
├── go.mod              # Go module definition
└── Makefile            # Build automation
```

## How It Works

[Implementation details to be added]

## Testing

[Testing instructions to be added]

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[License information to be added]
