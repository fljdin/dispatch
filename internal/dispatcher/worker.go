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

func NewWorker(memory *Memory, ctx context.Context) *Worker {
	return &Worker{
		memory:  memory,
		context: ctx,
	}
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
				w.handle(task)
			}
		}
	}
}

func (w *Worker) handle(t tasks.Task) {
	switch t.Status {
	case tasks.Waiting:
		w.memory.ForwardTask(t)

	case tasks.Interrupted:
		w.memory.results <- tasks.Result{
			ID:      t.ID,
			SubID:   t.SubID,
			Status:  tasks.Interrupted,
			Elapsed: 0,
		}

	case tasks.Ready:
		w.run(t)
	}
}

func (w *Worker) run(t tasks.Task) {
	if t.Command.From != "" {
		result, commands := t.Command.Generate()
		result.ID = t.ID
		result.SubID = t.SubID
		result.WorkerID = w.ID

		w.memory.results <- result

		for id, command := range commands {
			w.memory.AddTask(tasks.Task{
				ID:      t.ID,
				SubID:   id + 1,
				Command: command,
			})
		}
		return
	}

	result := t.Command.Run()
	result.ID = t.ID
	result.SubID = t.SubID
	result.WorkerID = w.ID

	w.memory.results <- result
}
