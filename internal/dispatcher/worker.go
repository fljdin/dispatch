package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/models"
)

type Worker struct {
	ID      int
	memory  *SharedMemory
	context context.Context
}

func (w *Worker) Start() {
	w.memory.StartWorker()
	defer w.memory.EndWorker()

	for {
		select {
		case <-w.context.Done():
			return
		case task := <-w.memory.tasks:
			status := w.getTaskStatus(task)

			if status == models.Ready {
				result := task.Command.Run()
				result.ID = task.ID
				result.QueryID = task.QueryID
				result.WorkerID = w.ID
				w.memory.results <- result
				continue
			}

			if status == models.Interrupted {
				w.memory.results <- models.TaskResult{
					ID:      task.ID,
					QueryID: task.QueryID,
					Status:  models.Interrupted,
					Elapsed: 0,
				}
				continue
			}

			w.memory.ForwardTask(task)
		}
	}
}

func (w *Worker) getTaskStatus(task models.Task) int {
	for _, id := range task.Depends {
		parentStatus := w.memory.GetStatus(id)

		if parentStatus >= models.Failed {
			return models.Interrupted
		} else if parentStatus < models.Succeeded {
			return models.Waiting
		}
	}

	return models.Ready
}
