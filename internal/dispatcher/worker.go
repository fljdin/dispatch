package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/task"
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

func (w *Worker) runTask(t *task.Task) {
	if t == nil {
		return
	}

	if t.Status == task.Ready {
		result := t.Command.Run()
		result.ID = t.ID
		result.QueryID = t.QueryID
		result.WorkerID = w.ID
		w.memory.results <- result
		return
	}

	if t.Status == task.Interrupted {
		w.memory.results <- task.TaskResult{
			ID:      t.ID,
			QueryID: t.QueryID,
			Status:  task.Interrupted,
			Elapsed: 0,
		}
		return
	}

	w.memory.ForwardTask(*t)
}
