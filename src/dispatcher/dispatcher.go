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
	plan       map[int]models.TaskResult
	readyTasks chan models.Task
	Results    chan models.TaskResult
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

	d.plan = make(map[int]models.TaskResult)
	d.readyTasks = make(chan models.Task, size)
	d.Results = make(chan models.TaskResult, size)

	// launch workers
	for i := 0; i < count; i++ {
		d.wgWorkers.Add(1)
		go d.worker(ctx)
	}

	return d
}

func (d *Dispatcher) Add(task models.Task) {
	var status = models.Waiting

	if len(task.Depends) == 0 {
		status = models.Ready
		d.wgTasks.Add(1)
		d.readyTasks <- task
	}

	d.plan[task.ID] = models.TaskResult{
		ID:     task.ID,
		Task:   &task,
		Status: status,
	}
}

func (d *Dispatcher) GetResult(ID int) models.TaskResult {
	return d.plan[ID]
}

func (d *Dispatcher) Wait() {
	// launch observer for tasks completion
	d.wgObserver.Add(1)
	go d.observer(d.context)

	d.wgTasks.Wait()    // wait until each task has been processed
	d.cancel()          // warm workers to stop theirs loop
	d.wgWorkers.Wait()  // wait until each worker has been stopped
	d.wgObserver.Wait() // wait for observer
}

func (d *Dispatcher) worker(ctx context.Context) {
	defer d.wgWorkers.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-d.readyTasks:
			d.Results <- task.Run(ctx)
			d.wgTasks.Done()
		}
	}
}

func (d *Dispatcher) observer(ctx context.Context) {
	defer d.wgObserver.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-d.Results:
			t := d.GetResult(result.ID)

			t.StartTime = result.StartTime
			t.EndTime = result.EndTime
			t.Elapsed = result.Elapsed
			t.Status = result.Status
			t.Output = result.Output
			t.Error = result.Error

			d.plan[t.ID] = t
			d.logger(t)
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
