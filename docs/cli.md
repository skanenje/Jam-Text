# CLI Package

Command-line interface for JamText text analysis operations.

## Commands
- `index` - Create searchable index
- `lookup` - Exact SimHash lookup
- `fuzzy` - Fuzzy SimHash lookup
- `hash` - Generate document hash
- `compare` - Compare two documents
- `moderate` - Content moderation

## Usage
```bash
jamtext -c <command> [options]

Options:
  -i string      Input file path
  -o string      Output file path
  -s int         Chunk size (default: 4096)
  -overlap int   Overlap size (default: 256)
  -threshold int Similarity threshold (default: 3)
```

## Examples
```bash
# Index creation
jamtext -c index -i book.txt -o book.idx -s 4096

# Fuzzy search
jamtext -c fuzzy -i book.idx -h <hash> -threshold 5
```