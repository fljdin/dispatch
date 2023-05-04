package dispatcher

import (
	"github.com/fljdin/dispatch/internal/models"
)

type Worker struct {
	ID         int
	dispatcher *Dispatcher
}

func (w *Worker) Start() {
	defer w.dispatcher.wgWorkers.Done()

	for {
		select {
		case <-w.dispatcher.context.Done():
			return
		case task := <-w.dispatcher.tasks:
			if len(task.Depends) == 0 {
				result := task.Run()
				result.WorkerID = w.ID
				w.dispatcher.results <- result
				continue
			}

			// verify if some dependencies have been completed
			var depends = []int{}
			var currentStatus = models.Waiting

			for _, id := range task.Depends {
				parentStatus := w.dispatcher.statuses.Load(id)

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
				w.dispatcher.results <- models.TaskResult{
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
			w.dispatcher.tasks <- task
		}
	}
}
