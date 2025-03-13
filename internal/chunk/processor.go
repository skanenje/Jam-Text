package chunk

import (
	"runtime"
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
