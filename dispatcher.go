package main

import (
	"context"
	"sync"
)

type Dispatcher struct {
	workerCount int
	wgWorkers   sync.WaitGroup
	wgTasks     sync.WaitGroup

	Tasks   chan Task
	Results chan TaskResult

	cancel func()
}

func NewDispatcher(ctx context.Context, count int, size int) *Dispatcher {
	ctx, cancel := context.WithCancel(ctx)

	d := &Dispatcher{
		workerCount: count,
		cancel:      cancel,
	}

	d.Tasks = make(chan Task, size)
	d.Results = make(chan TaskResult, size)

	for i := 0; i < d.workerCount; i++ {
		d.wgWorkers.Add(1)
		go d.worker(ctx)
	}

	return d
}

func (d *Dispatcher) Add(task Task) {
	d.wgTasks.Add(1)
	d.Tasks <- task
}

func (d *Dispatcher) Wait() {
	d.wgTasks.Wait()   // wait until each task has been proceed
	d.cancel()         // warm workers to stop their loop
	d.wgWorkers.Wait() // wait until each worker has been stopped
}

func (d *Dispatcher) worker(ctx context.Context) {
	defer d.wgWorkers.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-d.Tasks:
			d.Results <- task.Run(ctx)
			d.wgTasks.Done()
		}
	}
}
