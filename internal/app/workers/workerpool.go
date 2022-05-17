package workers

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type WorkerPool struct {
	workers int
	inputCh chan func(ctx context.Context) error
	done    chan struct{}
}

func New(ctx context.Context, workers int, buffer int) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		inputCh: make(chan func(ctx context.Context) error, buffer),
		done:    make(chan struct{}),
	}
}

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

func (wp *WorkerPool) Stop() {
	close(wp.done)
}

func (wp *WorkerPool) Push(task func(ctx context.Context) error) {
	wp.inputCh <- task
}
