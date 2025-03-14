# Chunk Package

Text segmentation and processing with context awareness.

## Types
```go
type Chunk struct {
    Content     string
    StartOffset int64
    Length      int
    IsComplete  bool
    Metadata    map[string]string
}

type ChunkOptions struct {
    ChunkSize        int    // Default: 4096
    OverlapSize      int    // Default: 256
    SplitOnBoundary  bool   // Default: true
    BoundaryChars    string // Default: ".!?\n"
    MaxChunkSize     int    // Default: 6144
    PreserveNewlines bool   // Default: true
}
```

## Usage
```go
// Create processor
processor := NewChunkProcessor(4, hyperplanes)
defer processor.Close()

// Process chunks
chunk := NewChunk("sample text", 0)
processor.ProcessChunk(chunk)

// Get results
for result := range processor.Results() {
    // Handle result
}
```

## Supported Formats
- Plain text (.txt)
- PDF documents (.pdf)
- Word documents (.docx)