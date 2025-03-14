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

func ReadChunk(filename string, position int64, chunkSize int) (chunk string, err error) {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".pdf":
		return readPDFChunk(filename, position, chunkSize)
	case ".docx":
		return readDocxChunk(filename, position, chunkSize)
	default:
		return readTextChunk(filename, position, chunkSize)
	}
}

func readTextChunk(filename string, position int64, chunkSize int) (chunk string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the content
	buffer := make([]byte, chunkSize)
	_, err = file.Seek(position, 0)
	if err != nil {
		return "", err
	}

	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Extract the chunk
	chunk = string(buffer[:bytesRead])

	return chunk, nil
}

func readPDFChunk(filename string, position int64, chunkSize int) (chunk string, err error) {
	// Check if pdftotext is installed
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", fmt.Errorf("pdftotext not found. Please install poppler-utils")
	}

	// Create a temporary file for the text output
	tmpFile, err := os.CreateTemp("", "pdf_*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Run pdftotext
	cmd := exec.Command("pdftotext", filename, tmpFile.Name())
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read the text file chunk
	return readTextChunk(tmpFile.Name(), position, chunkSize)
}

func readDocxChunk(filename string, position int64, chunkSize int) (chunk string, err error) {
	// Check if pandoc is installed
	if _, err := exec.LookPath("pandoc"); err != nil {
		return "", fmt.Errorf("pandoc not found. Please install pandoc")
	}

	// Create a temporary file for the text output
	tmpFile, err := os.CreateTemp("", "docx_*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Run pandoc to convert docx to text
	cmd := exec.Command("pandoc", "-f", "docx", "-t", "plain", filename, "-o", tmpFile.Name())
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error converting docx: %w", err)
	}

	// Read the text file chunk
	return readTextChunk(tmpFile.Name(), position, chunkSize)
}
