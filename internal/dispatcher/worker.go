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
			task, result := w.verifyStatus(task)

			if result.Status == models.Ready {
				result := task.Command.Run()
				result.ID = task.ID
				result.QueryID = task.QueryID
				result.WorkerID = w.ID
				w.memory.results <- result
				continue
			}

			if result.Status == models.Interrupted {
				w.memory.results <- result
				continue
			}

			w.memory.ForwardTask(task)
		}
	}
}

func (w *Worker) verifyStatus(task models.Task) (models.Task, models.TaskResult) {
	var depends = []int{}
	var result models.TaskResult = models.TaskResult{
		Status: models.Waiting,
	}

	for _, id := range task.Depends {
		parentStatus := w.memory.GetStatus(id)

		if parentStatus < models.Succeeded {
			depends = append(depends, id)
			continue
		}

		if parentStatus >= models.Failed {
			return task, models.TaskResult{
				ID:       task.ID,
				QueryID:  task.QueryID,
				WorkerID: w.ID,
				Status:   models.Interrupted,
				Elapsed:  0,
			}
		}
	}

	task.Depends = depends

	if len(depends) == 0 {
		result.Status = models.Ready
	}

	return task, result
}
