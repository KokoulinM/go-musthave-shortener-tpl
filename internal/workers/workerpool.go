// Package workers receive work on the jobs channel and send the corresponding results on results.
package workers

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type WorkerPool struct {
	// workers -number of workers
	workers int
	// inputCh - these channel will receive work
	inputCh chan func(ctx context.Context) error
	// done - channel for stopping the work of the worker
	done chan struct{}
}

// New is the worker constructor
func New(ctx context.Context, workers int, buffer int) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		inputCh: make(chan func(ctx context.Context) error, buffer),
		done:    make(chan struct{}),
	}
}

// Run is the method to start the worker
func (wp *WorkerPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}

	for i := 0; i < wp.workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("Worker #%v start \n", i)
		outer:
			for {
				select {
				case f := <-wp.inputCh:
					err := f(ctx)
					if err != nil {
						fmt.Printf("Error on worker #%v: %v\n", i, err.Error())
					}
				case <-wp.done:
					break outer
				}
			}
			log.Printf("Worker #%v close\n", i)
		}(i)
	}
	wg.Wait()
	close(wp.inputCh)
}

// Stop is the method to stop the worker
func (wp *WorkerPool) Stop() {
	close(wp.done)
}

// Push is the method to push into the inputCh
func (wp *WorkerPool) Push(task func(ctx context.Context) error) {
	wp.inputCh <- task
}
