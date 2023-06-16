package dispatcher

import (
	"context"
	"sync"

	"github.com/fljdin/dispatch/internal/logger"
	"github.com/fljdin/dispatch/internal/models"
	"github.com/sourcegraph/conc/pool"
)

type Dispatcher struct {
	workers  int
	tasks    chan models.Task
	results  chan models.TaskResult
	statuses models.StatusMap

	wgTasks sync.WaitGroup
	context context.Context
	cancel  func()
}

func (d *Dispatcher) AddTask(task models.Task) {
	d.tasks <- task
	d.wgTasks.Add(1)
}

func (d *Dispatcher) GetStatus(ID int) int {
	return d.statuses.Load(ID)
}

func (d *Dispatcher) Wait() {
	workers := pool.New().
		WithContext(d.context).
		WithMaxGoroutines(d.workers)

	// WARNING:
	// as run() can postpone task, we cannot close chan
	// then, conc.pool will loop forever even if we use
	// a WaitGroup on task completion
	for task := range d.tasks {
		task := task
		workers.Go(func(ctx context.Context) error {
			d.run(task)
			return nil
		})
	}

	observer := pool.New().WithMaxGoroutines(1)
	for result := range d.results {
		result := result
		observer.Go(func() {
			logger := logger.Console{}
			logger.Render(result)
			d.wgTasks.Done()
		})
	}

	d.wgTasks.Wait()
	d.cancel()
	workers.Wait()

	close(d.results)
	observer.Wait()
}

func (d *Dispatcher) run(task models.Task) {
	status := d.getTaskStatus(task)

	if status == models.Ready {
		result := task.Command.Run()
		result.ID = task.ID
		result.QueryID = task.QueryID

		d.statuses.Store(task.ID, result.Status)

		d.results <- result
		return
	}

	if status == models.Interrupted {
		d.results <- models.TaskResult{
			ID:      task.ID,
			QueryID: task.QueryID,
			Status:  models.Interrupted,
			Elapsed: 0,
		}
		return
	}

	// postpone
	d.AddTask(task)
}

func (d *Dispatcher) getTaskStatus(task models.Task) int {
	for _, id := range task.Depends {
		parentStatus := d.GetStatus(id)

		if parentStatus >= models.Failed {
			return models.Interrupted
		} else if parentStatus < models.Succeeded {
			return models.Waiting
		}
	}

	return models.Ready
}
