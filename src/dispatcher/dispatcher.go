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
	completed  sync.Map
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
	d.completed = sync.Map{}

	// launch observer
	d.wgObserver.Add(1)
	go d.observer(d.context)

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

func (d *Dispatcher) GetResult(ID int) (models.TaskResult, bool) {
	if result, ok := d.completed.Load(ID); ok {
		return result.(models.TaskResult), ok
	}
	return models.TaskResult{}, false
}

func (d *Dispatcher) Wait() {
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
				continue
			}

			// verify if some dependencies have been completed
			var depends = []int{}
			var status = models.Waiting

			for _, id := range task.Depends {
				dependency, exists := d.completed.Load(id)

				if !exists {
					// dependency has not been completed yet
					depends = append(depends, id)
					continue
				}

				if dependency.(models.TaskResult).Status >= models.Failed {
					status = models.Interrupted
				}
			}

			// current task is interrupted and won't be launched
			if status == models.Interrupted {
				d.results <- models.TaskResult{
					ID:      task.ID,
					Status:  status,
					Elapsed: 0,
				}
				continue
			}

			// forward task to another worker
			task.Depends = depends
			d.tasks <- task
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
			d.completed.Store(result.ID, result)
			d.logger(result)
			d.wgTasks.Done()
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
