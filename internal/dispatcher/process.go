package dispatcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Process struct {
	ID      int
	memory  *Memory
	context context.Context
}

func (p *Process) Start() {
	p.memory.StartProcess()
	defer p.memory.EndProcess()

	for {
		select {
		case <-p.context.Done():
			return
		default:
			if task, ok := p.memory.queue.Pop(); ok {
				p.handle(task)
			}
		}
	}
}

func (p *Process) handle(t tasks.Task) {
	switch t.Status {
	case tasks.Waiting:
		p.memory.ForwardTask(t)

	case tasks.Interrupted:
		p.memory.results <- tasks.Result{
			ID:      t.ID,
			SubID:   t.SubID,
			Name:    t.Name,
			Action:  t.Action.String(),
			Status:  tasks.Interrupted,
			Elapsed: 0,
		}

	case tasks.Ready:
		p.run(t)
	}
}

func (p *Process) run(t tasks.Task) {
	report, commands := t.Action.Run()

	for id, command := range commands {
		p.memory.AddTask(tasks.Task{
			ID:     t.ID,
			SubID:  id + 1,
			Action: command,
			Name:   t.Name,
		})
	}

	logAttrs := []any{
		"status", tasks.StatusTypes[report.Status],
		"name", t.Name,
		"start", report.StartTime.Format(time.DateTime),
		"elapsed", report.Elapsed.Round(time.Millisecond),
	}

	if !tasks.IsSucceeded(report.Status) {
		slog.Error(t.Code(), logAttrs...)
	} else {
		slog.Info(t.Code(), logAttrs...)
	}
	slog.Info(t.Code(), "action", t.Action.String())

	if len(report.Error) > 0 {
		msg := fmt.Sprintf("%s Error:\n%s", t.Code(), report.Error)
		slog.Error(msg)
	}

	if len(report.Output) > 0 {
		msg := fmt.Sprintf("%s Output:\n%s", t.Code(), report.Output)
		slog.Info(msg)
	}

	p.memory.results <- tasks.Result{
		ID:        t.ID,
		SubID:     t.SubID,
		Name:      t.Name,
		Action:    t.Action.String(),
		Status:    report.Status,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
		Elapsed:   report.Elapsed,
		Output:    report.Output,
		Error:     report.Error,
	}
}
