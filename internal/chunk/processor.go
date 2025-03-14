package chunk

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"time"
	"unicode/utf8"

	"jamtext/internal/index"
	"jamtext/internal/simhash"
)

// ChunkProcessor handles the processing of file chunks
type ChunkProcessor struct {
	pool        *WorkerPool
	resultChan  chan ProcessResult
	vectorizer  simhash.Vectorizer
	hyperplanes [][]float64
}

// ProcessResult represents the result of processing a chunk
type ProcessResult struct {
	Hash  simhash.SimHash
	Pos   int64
	Error error
}

// NewChunkProcessor creates a new chunk processor
func NewChunkProcessor(numWorkers int, hyperplanes [][]float64) *ChunkProcessor {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	return &ChunkProcessor{
		pool:        NewWorkerPool(numWorkers),
		resultChan:  make(chan ProcessResult, numWorkers*2),
		vectorizer:  simhash.NewFrequencyVectorizer(simhash.VectorDimensions),
		hyperplanes: hyperplanes,
	}
}

// ProcessChunk handles the processing of a single chunk
func (cp *ChunkProcessor) ProcessChunk(chunk Chunk) {
	cp.pool.Submit(func() {
		hash := simhash.CalculateWithVectorizer(chunk.Content, cp.hyperplanes, cp.vectorizer)
		cp.resultChan <- ProcessResult{
			Hash: hash,
			Pos:  chunk.StartOffset,
		}
	})
}

// Close shuts down the chunk processor
func (cp *ChunkProcessor) Close() {
	cp.pool.Close()
	close(cp.resultChan)
}

// Results returns the channel for receiving processing results
func (cp *ChunkProcessor) Results() <-chan ProcessResult {
	return cp.resultChan
}

// isValidUTF8 checks if a byte slice is valid UTF-8
func isValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// findBoundary finds a good boundary to split text
func findBoundary(text []byte, preferredPos int, boundaryChars string) int {
	if preferredPos >= len(text) {
		return len(text)
	}

	// Check backward from preferred position
	for i := preferredPos; i > max(0, preferredPos-100); i-- {
		for _, c := range boundaryChars {
			if i < len(text) && text[i] == byte(c) {
				return i + 1
			}
		}
	}

	// If no good boundary found, use the preferred position
	return preferredPos
}

// ProcessFile chunks a file and builds an index with advanced options
func ProcessFile(filename string, opts ChunkOptions, hyperplanes [][]float64, indexDir string) (*index.Index, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	idx := index.New(filename, opts.ChunkSize, hyperplanes, indexDir)

	// Create chunk processor
	processor := NewChunkProcessor(runtime.NumCPU(), hyperplanes)

	// Start result consumer
	resultsDone := make(chan struct{})
	go func() {
		defer close(resultsDone)
		count := 0
		for result := range processor.Results() {
			if result.Error != nil {
				opts.Logger.Printf("Error processing chunk: %v", result.Error)
				continue
			}

			if opts.Verbose {
				opts.Logger.Printf("Chunk %d: offset=%d, hash=%016x",
					count, result.Pos, result.Hash)
			}

			// Log every hash
			if true { // Changed from if count%100 == 0
				opts.Logger.Printf("Hash: %016x at position %d",
					result.Hash, result.Pos)
			}

			idx.Add(result.Hash, result.Pos)
			count++
		}
	}()

	reader := bufio.NewReader(file)
	buffer := make([]byte, opts.ChunkSize)
	offset := int64(0)
	overlap := make([]byte, 0, opts.OverlapSize)

	// Process chunks
	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			processor.Close() // Close processor on error
			return nil, err
		}
		if bytesRead > 0 {
			chunkData := buffer[:bytesRead]

			// Handle overlap from previous chunk
			if len(overlap) > 0 {
				combinedChunk := make([]byte, len(overlap)+bytesRead)
				copy(combinedChunk, overlap)
				copy(combinedChunk[len(overlap):], chunkData)
				chunkData = combinedChunk
			}

			// Fix UTF-8 encoding issues at chunk boundaries
			if !isValidUTF8(chunkData) {
				// Try to find a valid UTF-8 boundary
				validLen := 0
				for i := 0; i < len(chunkData); i++ {
					if isValidUTF8(chunkData[:i+1]) {
						validLen = i + 1
					}
				}
				if validLen > 0 {
					chunkData = chunkData[:validLen]
				}
			}

			// Find a good boundary to split if needed
			chunkSize := len(chunkData)
			splitPos := chunkSize
			if opts.SplitOnBoundary && chunkSize > opts.ChunkSize/2 {
				splitPos = findBoundary(chunkData, opts.ChunkSize, opts.BoundaryChars)
			}

			// Create the chunk
			chunk := Chunk{
				Content:     string(chunkData[:splitPos]),
				StartOffset: offset,
				Length:      splitPos,
				IsComplete:  err == io.EOF,
				Metadata: map[string]string{
					"timestamp": time.Now().Format(time.RFC3339),
				},
			}

			processor.ProcessChunk(chunk)

			// Prepare overlap for next chunk
			if splitPos < chunkSize && opts.OverlapSize > 0 {
				overlapStart := max(0, splitPos-opts.OverlapSize)
				overlap = make([]byte, splitPos-overlapStart)
				copy(overlap, chunkData[overlapStart:splitPos])
			} else {
				overlap = nil
			}

			offset += int64(splitPos)
		}

		if err == io.EOF {
			break
		}
	}

	// Close the processor BEFORE waiting for results
	processor.Close()

	// Wait for all results to be processed
	<-resultsDone

	return idx, nil
}
