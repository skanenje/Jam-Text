# Jam-Text

Thia is an implementation of a text indexer using SimHash fingerprints based on vector similarity with random hyperlanes. The implementation is written in Go and provides an efficient way to index and search large text files.

## The text indexer works by:
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

## How It Works

## Testing

## Contributors

## LICENSE



