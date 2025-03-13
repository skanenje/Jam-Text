package chunk

import (
	"context"
	"sync"
)

// WorkerPool manages a pool of worker goroutines
type WorkerPool struct {
	workers    int
	tasks      chan func()
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &WorkerPool{
		workers:    workers,
		tasks:      make(chan func(), workers*3), // Buffer channel for smoother operation
		ctx:        ctx,
		cancelFunc: cancel,
	}

	pool.start()
	return pool
}

// start initializes the worker goroutines
func (p *WorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					task()
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

// Submit adds a task to the worker pool
func (p *WorkerPool) Submit(task func()) {
	select {
	case p.tasks <- task:
		// Task submitted
	case <-p.ctx.Done():
		// Pool has been stopped
	}
}

// Close stops the worker pool
func (p *WorkerPool) Close() {
	p.cancelFunc()
	close(p.tasks)
	p.wg.Wait()
}
