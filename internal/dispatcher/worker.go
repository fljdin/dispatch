package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/queue"
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
	case queue.Waiting:
		w.memory.ForwardTask(t)

	case queue.Interrupted:
		w.memory.results <- tasks.Result{
			ID:      t.ID,
			SubID:   t.SubID,
			Status:  queue.Interrupted,
			Elapsed: 0,
		}

	case queue.Ready:
		w.run(t)
	}
}

func (w *Worker) run(t tasks.Task) {
	report, commands := t.Action.Run()

	for id, command := range commands {
		w.memory.AddTask(tasks.Task{
			ID:     t.ID,
			SubID:  id + 1,
			Action: command,
		})
	}

	w.memory.results <- tasks.Result{
		ID:        t.ID,
		SubID:     t.SubID,
		WorkerID:  w.ID,
		Status:    report.Status,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
		Elapsed:   report.Elapsed,
		Output:    report.Output,
		Error:     report.Error,
	}
}
