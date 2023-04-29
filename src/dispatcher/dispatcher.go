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
	for i := 1; i <= count; i++ {
		worker := &Worker{
			ID:         i,
			dispatcher: d,
		}
		go worker.Start()
		d.wgWorkers.Add(1)
	}

	return d
}

func (d *Dispatcher) Add(task models.Task) {
	d.tasks <- task
	d.wgTasks.Add(1)
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
		"Worker %d completed Task %d (success: %t, elapsed: %s)\n",
		result.WorkerID,
		result.ID,
		(result.Status == models.Succeeded),
		result.Elapsed.Round(time.Millisecond),
	)
}
