package chunk

import (
	"io"
	"log"
	"os"
)

// Chunk represents a text chunk
type Chunk struct {
	StartOffset int64
	Content     string
	IsComplete  bool // Indicates if this is a complete text unit
	Metadata    map[string]string
}

// ChunkOptions represents options for chunk processing
type ChunkOptions struct {
	ChunkSize        int
	OverlapSize      int
	SplitOnBoundary  bool
	BoundaryChars    string
	MaxChunkSize     int
	PreserveNewlines bool
	Logger           *log.Logger
	Verbose          bool
}

// DefaultChunkOptions returns default chunking options
func DefaultChunkOptions() ChunkOptions {
	return ChunkOptions{
		ChunkSize:        4096,
		OverlapSize:      256,
		SplitOnBoundary:  true,
		BoundaryChars:    ".!?\n",
		MaxChunkSize:     6144,
		PreserveNewlines: true,
	}
}

// ReadChunk reads content at a specific position with context
func ReadChunk(filename string, position int64, chunkSize int, contextBefore, contextAfter int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Calculate start position with context
	startPos := position - int64(contextBefore)
	if startPos < 0 {
		startPos = 0
	}

	// Calculate total size with context
	totalSize := chunkSize + contextBefore + contextAfter
	buffer := make([]byte, totalSize)

	// Seek to the start position
	if _, err := file.Seek(startPos, 0); err != nil {
		return "", err
	}

	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Adjust for the actual amount read
	actualContent := buffer[:bytesRead]

	// Ensure we have valid UTF-8
	if !isValidUTF8(actualContent) {
		// Try to find a valid UTF-8 boundary
		validLen := 0
		for i := 0; i < len(actualContent); i++ {
			if isValidUTF8(actualContent[:i+1]) {
				validLen = i + 1
			}
		}
		if validLen > 0 {
			actualContent = actualContent[:validLen]
		}
	}

	return string(actualContent), nil
}
