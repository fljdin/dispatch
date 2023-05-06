package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/models"
)

type Worker struct {
	ID  int
	mem *WorkerMem
	ctx context.Context
}

func (w *Worker) Start() {
	w.mem.wgWorkers.Add(1)
	defer w.mem.wgWorkers.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case task := <-w.mem.tasks:
			if len(task.Depends) == 0 {
				result := task.Run()
				result.WorkerID = w.ID
				w.mem.results <- result
				continue
			}

			// verify if some dependencies have been completed
			var depends = []int{}
			var currentStatus = models.Waiting

			for _, id := range task.Depends {
				parentStatus := w.mem.statuses.Load(id)

				if parentStatus == models.Waiting {
					depends = append(depends, id)
					continue
				}

				if parentStatus >= models.Failed {
					currentStatus = models.Interrupted
				}
			}

			// current task is interrupted and won't be launched
			if currentStatus == models.Interrupted {
				w.mem.results <- models.TaskResult{
					ID:       task.ID,
					QueryID:  task.QueryID,
					WorkerID: w.ID,
					Status:   currentStatus,
					Elapsed:  0,
				}
				continue
			}

			// forward task to another worker
			task.Depends = depends
			w.mem.tasks <- task
		}
	}
}
