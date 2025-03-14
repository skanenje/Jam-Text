# Chunk Package Documentation

## Overview
The chunk package handles text segmentation and processing for JamText, supporting multiple file formats and providing efficient chunk management with context awareness.

## Core Types

### Chunk
```go
type Chunk struct {
    Content     string
    StartOffset int64
    Length      int
    IsComplete  bool
    Metadata    map[string]string
}
```

### ChunkOptions
```go
type ChunkOptions struct {
    ChunkSize        int    // Default: 4096
    OverlapSize      int    // Default: 256
    SplitOnBoundary  bool   // Default: true
    BoundaryChars    string // Default: ".!?\n"
    MaxChunkSize     int    // Default: 6144
    PreserveNewlines bool   // Default: true
    Logger           Logger
    Verbose          bool
}
```

## Functions

### NewChunk
```go
func NewChunk(content string, startOffset int64) Chunk
```
Creates a new chunk with initialized fields.
- **Parameters:**
  - `content`: The text content of the chunk
  - `startOffset`: Starting position in the source file
- **Returns:** Initialized Chunk structure

### ReadChunk
```go
func ReadChunk(filename string, position int64, chunkSize int, contextBefore, contextAfter int) (chunk string, contextBeforeStr string, contextAfterStr string, err error)
```
Reads a chunk from a file with context.
- **Parameters:**
  - `filename`: Path to the file
  - `position`: Starting position
  - `chunkSize`: Size of chunk to read
  - `contextBefore`: Number of bytes before chunk
  - `contextAfter`: Number of bytes after chunk
- **Supported Formats:**
  - Plain text (.txt)
  - PDF documents (.pdf)
  - Word documents (.docx)
- **Returns:**
  - The chunk content
  - Context before the chunk
  - Context after the chunk
  - Any error encountered

### ProcessFile
```go
func ProcessFile(filename string, opts ChunkOptions, hyperplanes [][]float64, indexDir string) (*index.Index, error)
```
Processes a file into chunks and builds an index.
- **Parameters:**
  - `filename`: Path to input file
  - `opts`: Chunking configuration
  - `hyperplanes`: SimHash hyperplanes
  - `indexDir`: Directory for index storage
- **Returns:**
  - Constructed index
  - Any error encountered

### ChunkProcessor Methods

#### NewChunkProcessor
```go
func NewChunkProcessor(numWorkers int, hyperplanes [][]float64) *ChunkProcessor
```
Creates a new chunk processor with worker pool.
- **Parameters:**
  - `numWorkers`: Number of concurrent workers
  - `hyperplanes`: SimHash hyperplanes
- **Returns:** Initialized ChunkProcessor

#### ProcessChunk
```go
func (cp *ChunkProcessor) ProcessChunk(chunk Chunk)
```
Processes a single chunk asynchronously.
- **Parameters:**
  - `chunk`: Chunk to process

#### Close
```go
func (cp *ChunkProcessor) Close()
```
Shuts down the chunk processor and its worker pool.

#### Results
```go
func (cp *ChunkProcessor) Results() <-chan ProcessResult
```
Returns channel for receiving processing results.

## Helper Functions

### findBoundary
```go
func findBoundary(text []byte, preferredPos int, boundaryChars string) int
```
Finds appropriate text boundary for splitting.
- **Parameters:**
  - `text`: Input text
  - `preferredPos`: Preferred split position
  - `boundaryChars`: Valid boundary characters
- **Returns:** Actual split position

### isValidUTF8
```go
func isValidUTF8(data []byte) bool
```
Checks if byte slice is valid UTF-8.
- **Parameters:**
  - `data`: Bytes to check
- **Returns:** true if valid UTF-8

## Usage Examples

### Basic Chunk Processing
```go
opts := ChunkOptions{
    ChunkSize:       4096,
    OverlapSize:     256,
    SplitOnBoundary: true,
    BoundaryChars:   ".!?\n",
}

idx, err := ProcessFile("document.txt", opts, hyperplanes, "index_dir")
if err != nil {
    log.Fatal(err)
}
```

### Reading with Context
```go
chunk, before, after, err := ReadChunk("document.pdf", 1000, 4096, 100, 100)
if err != nil {
    log.Fatal(err)
}
```

### Custom Chunk Processing
```go
processor := NewChunkProcessor(4, hyperplanes)
defer processor.Close()

chunk := NewChunk("sample text", 0)
processor.ProcessChunk(chunk)

for result := range processor.Results() {
    // Handle result
}
```

## Dependencies
- `pdftotext` (poppler-utils) for PDF support
- `pandoc` for DOCX support