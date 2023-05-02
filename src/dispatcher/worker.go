package dispatcher

import (
	"github.com/fljdin/dispatch/src/models"
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
			var status = models.Waiting

			for _, id := range task.Depends {
				completed, exists := w.dispatcher.completed.Load(id)

				if !exists {
					// dependency has not been completed yet
					depends = append(depends, id)
					continue
				}

				if completed.(int) >= models.Failed {
					status = models.Interrupted
				}
			}

			// current task is interrupted and won't be launched
			if status == models.Interrupted {
				w.dispatcher.results <- models.TaskResult{
					ID:       task.ID,
					QueryID:  task.QueryID,
					WorkerID: w.ID,
					Status:   status,
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
