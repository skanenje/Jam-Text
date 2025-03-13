package chunk

import (
	"testing"
)

func TestDefaultChunkOptions(t *testing.T) {
	opts := DefaultChunkOptions()

	if opts.ChunkSize != 4096 {
		t.Errorf("Expected ChunkSize 4096, got %d", opts.ChunkSize)
	}
	if opts.OverlapSize != 256 {
		t.Errorf("Expected OverlapSie 256, got %d", opts.OverlapSize)
	}

	if !opts.SplitOnBoundary {
		t.Errorf("Expected SplitOnBoundary to be true")
	}
}
<<<<<<< HEAD

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(4)
	defer pool.Close()

	results := make(chan int, 10)
	for i := 0; i < 10; i++ {
		i := i
		pool.Submit(func() {
			results <- i * 2
		})
	}

	// Collecting results
	sum := 0
	for i := 0; i < 10; i++ {
		sum += <-results
	}

	expected := 90 // sum of 0*2 through 9*2
	if sum != expected {
		t.Errorf("Expected sum %d, got %d", expected, sum)
	}
}
=======
>>>>>>> def53ac (test(chunk): add basic test for DefaultChunkOptions)
