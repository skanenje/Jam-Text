package chunk

import (
	"fmt"
	"io"
	"os"
	"strings"
	"path/filepath"
	"os/exec"
)

// Chunk represents a section of text with its metadata
type Chunk struct {
	Content     string
	StartOffset int64
	Length      int
	IsComplete  bool
	Metadata    map[string]string
}

// NewChunk creates a new chunk with initialized fields
func NewChunk(content string, startOffset int64) Chunk {
	return Chunk{
		Content:     content,
		StartOffset: startOffset,
		Length:      len(content),
		IsComplete:  true,
		Metadata:    make(map[string]string),
	}
}

// ChunkOptions defines configuration for chunk processing
type ChunkOptions struct {
	ChunkSize        int
	OverlapSize      int
	SplitOnBoundary  bool
	BoundaryChars    string
	MaxChunkSize     int
	PreserveNewlines bool
	Logger           Logger
	Verbose          bool
}

// Logger interface for logging operations
type Logger interface {
	Printf(format string, v ...interface{})
}

func ReadChunk(filename string, position int64, chunkSize int, contextBefore, contextAfter int) (chunk string, contextBeforeStr string, contextAfterStr string, err error) {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".pdf":
		return readPDFChunk(filename, position, chunkSize, contextBefore, contextAfter)
	case ".docx":
		return readDocxChunk(filename, position, chunkSize, contextBefore, contextAfter)
	default:
		return readTextChunk(filename, position, chunkSize, contextBefore, contextAfter)
	}
}

func readTextChunk(filename string, position int64, chunkSize int, contextBefore, contextAfter int) (chunk string, contextBeforeStr string, contextAfterStr string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	// Read the content
	buffer := make([]byte, contextBefore+chunkSize+contextAfter)
	_, err = file.Seek(position-int64(contextBefore), 0)
	if err != nil {
		return "", "", "", err
	}

	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", "", "", err
	}

	// Calculate boundaries
	contextBeforeStart := 0
	chunkStart := contextBefore
	chunkEnd := chunkStart + chunkSize
	if chunkEnd > bytesRead {
		chunkEnd = bytesRead
	}

	// Extract the chunks
	if contextBefore > 0 {
		contextBeforeStr = string(buffer[contextBeforeStart:chunkStart])
	}
	chunk = string(buffer[chunkStart:chunkEnd])
	if chunkEnd < bytesRead {
		contextAfterStr = string(buffer[chunkEnd:bytesRead])
	}

	return chunk, contextBeforeStr, contextAfterStr, nil
}

func readPDFChunk(filename string, position int64, chunkSize int, contextBefore, contextAfter int) (chunk string, contextBeforeStr string, contextAfterStr string, err error) {
	// Check if pdftotext is installed
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", "", "", fmt.Errorf("pdftotext not found. Please install poppler-utils")
	}

	// Create a temporary file for the text output
	tmpFile, err := os.CreateTemp("", "pdf_*.txt")
	if err != nil {
		return "", "", "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Run pdftotext
	cmd := exec.Command("pdftotext", filename, tmpFile.Name())
	if err := cmd.Run(); err != nil {
		return "", "", "", err
	}

	// Read the text file chunk
	return readTextChunk(tmpFile.Name(), position, chunkSize, contextBefore, contextAfter)
}

func readDocxChunk(filename string, position int64, chunkSize int, contextBefore, contextAfter int) (chunk string, contextBeforeStr string, contextAfterStr string, err error) {
	// Check if pandoc is installed
	if _, err := exec.LookPath("pandoc"); err != nil {
		return "", "", "", fmt.Errorf("pandoc not found. Please install pandoc")
	}

	// Create a temporary file for the text output
	tmpFile, err := os.CreateTemp("", "docx_*.txt")
	if err != nil {
		return "", "", "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Run pandoc to convert docx to text
	cmd := exec.Command("pandoc", "-f", "docx", "-t", "plain", filename, "-o", tmpFile.Name())
	if err := cmd.Run(); err != nil {
		return "", "", "", fmt.Errorf("error converting docx: %w", err)
	}

	// Read the text file chunk
	return readTextChunk(tmpFile.Name(), position, chunkSize, contextBefore, contextAfter)
}
