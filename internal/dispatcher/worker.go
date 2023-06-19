package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/models"
)

type Worker struct {
	ID      int
	memory  *Memory
	context context.Context
}

func (w *Worker) Start() {
	w.memory.StartWorker()
	defer w.memory.EndWorker()

	for {
		select {
		case <-w.context.Done():
			return
		default:
			task := w.memory.queue.Pop()
			w.runTask(task)
		}
	}
}

func (w *Worker) runTask(task *models.Task) {
	if task == nil {
		return
	}

	if task.Status == models.Ready {
		result := task.Command.Run()
		result.ID = task.ID
		result.QueryID = task.QueryID
		result.WorkerID = w.ID
		w.memory.results <- result
		return
	}

	if task.Status == models.Interrupted {
		w.memory.results <- models.TaskResult{
			ID:      task.ID,
			QueryID: task.QueryID,
			Status:  models.Interrupted,
			Elapsed: 0,
		}
		return
	}

	w.memory.ForwardTask(*task)
}
