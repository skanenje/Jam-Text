package chunk

import "log"

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
