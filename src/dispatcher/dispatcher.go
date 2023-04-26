package dispatcher

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/fljdin/dispatch/src/models"
)

type Dispatcher struct {
	context    context.Context
	tasks      chan models.Task
	results    chan models.TaskResult
	completed  map[int]models.TaskResult
	wgTasks    sync.WaitGroup
	wgWorkers  sync.WaitGroup
	wgObserver sync.WaitGroup
	cancel     func()
}

func NewDispatcher(ctx context.Context, count int, size int) *Dispatcher {
	ctx, cancel := context.WithCancel(ctx)

	d := &Dispatcher{
		context: ctx,
		cancel:  cancel,
	}

	d.tasks = make(chan models.Task, size)
	d.results = make(chan models.TaskResult, size)
	d.completed = make(map[int]models.TaskResult)

	// launch workers
	for i := 0; i < count; i++ {
		d.wgWorkers.Add(1)
		go d.worker(ctx)
	}

	return d
}

func (d *Dispatcher) Add(task models.Task) {
	d.wgTasks.Add(1)
	d.tasks <- task
}

func (d *Dispatcher) GetResult(ID int) models.TaskResult {
	return d.completed[ID]
}

func (d *Dispatcher) Wait() {
	// launch observer
	d.wgObserver.Add(1)
	go d.observer(d.context)

	d.wgTasks.Wait()    // wait until each task has been processed
	d.cancel()          // warm workers to stop theirs loop
	d.wgWorkers.Wait()  // wait until each worker has been stopped
	d.wgObserver.Wait() // wait until observer has been stopped
}

func (d *Dispatcher) worker(ctx context.Context) {
	defer d.wgWorkers.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-d.tasks:
			if len(task.Depends) == 0 {
				d.results <- task.Run(ctx)
				d.wgTasks.Done()
			} else {
				// update depends if needed
				var depends = []int{}
				var status = models.Waiting

				for _, id := range task.Depends {
					result, exists := d.completed[id]

					if !exists {
						depends = append(depends, id)
					} else if result.Status >= models.Failed {
						status = models.Interrupted
					}
				}

				if status != models.Interrupted {
					task.Depends = depends
					d.tasks <- task
				} else {
					d.results <- models.TaskResult{
						ID:     task.ID,
						Status: models.Interrupted,
					}
					d.wgTasks.Done()
				}
			}
		}
	}
}

func (d *Dispatcher) observer(ctx context.Context) {
	defer d.wgObserver.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-d.results:
			d.completed[result.ID] = result
			d.logger(result)
		}
	}
}

func (d *Dispatcher) logger(result models.TaskResult) {
	log.Printf(
		"Task %d completed (success: %t, elapsed: %s)\n",
		result.ID,
		(result.Status == models.Succeeded),
		result.Elapsed.Round(time.Millisecond),
	)
}
