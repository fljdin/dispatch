package main

import (
	"context"
	"log"
	"sync"
	"time"
)

type Dispatcher struct {
	context   context.Context
	Tasks     chan Task
	Results   chan TaskResult
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	wgLogger  sync.WaitGroup
	cancel    func()
}

func NewDispatcher(ctx context.Context, count int, size int) *Dispatcher {
	ctx, cancel := context.WithCancel(ctx)

	d := &Dispatcher{
		context: ctx,
		cancel:  cancel,
	}

	d.Tasks = make(chan Task, size)
	d.Results = make(chan TaskResult, size)

	// launch workers
	for i := 0; i < count; i++ {
		d.wgWorkers.Add(1)
		go d.worker(ctx)
	}

	return d
}

func (d *Dispatcher) Add(task Task) {
	d.wgTasks.Add(1)
	d.Tasks <- task
}

// launch logger on demand
func (d *Dispatcher) Log() {
	d.wgLogger.Add(1)
	go d.logger(d.context)
}

func (d *Dispatcher) Wait() {
	d.wgTasks.Wait()   // wait until each task has been processed
	d.cancel()         // warm workers to stop theirs loop
	d.wgWorkers.Wait() // wait until each worker has been stopped
	d.wgLogger.Wait()  // wait for logger completion
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

func (d *Dispatcher) logger(ctx context.Context) {
	defer d.wgLogger.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-d.Results:
			log.Printf(
				"Task %d completed (success: %t, elapsed: %s)\n",
				result.ID,
				(result.Status == Succeeded),
				result.Elapsed.Round(time.Millisecond),
			)
		}
	}
}
