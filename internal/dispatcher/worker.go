package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/tasks"
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
			if task, ok := w.memory.queue.Pop(); ok {
				w.run(task)
			}
		}
	}
}

func (w *Worker) run(t tasks.Task) {
	if t.Status == tasks.Ready {
		result := t.Command.Run()
		result.ID = t.ID
		result.QueryID = t.QueryID
		result.WorkerID = w.ID
		w.memory.results <- result
		return
	}

	if t.Status == tasks.Interrupted {
		w.memory.results <- tasks.Result{
			ID:      t.ID,
			QueryID: t.QueryID,
			Status:  tasks.Interrupted,
			Elapsed: 0,
		}
		return
	}

	w.memory.ForwardTask(t)
}
