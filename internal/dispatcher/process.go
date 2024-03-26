package dispatcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fljdin/dispatch/internal/status"
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
		case task := <-p.memory.tasks:
			if task.Status == status.Interrupted {
				p.interrupt(task)
				continue
			}
			p.run(task)
		}
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
		"status", report.Status,
		"name", t.Name,
		"start", report.StartTime.Format(time.DateTime),
		"elapsed", report.Elapsed.Round(time.Millisecond),
	}

	if report.Status.IsFailed() {
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

	p.memory.results <- Result{
		ID:     t.ID,
		SubID:  t.SubID,
		Status: report.Status,
	}
}

func (p *Process) interrupt(t tasks.Task) {
	logAttrs := []any{
		"status", status.Interrupted,
		"name", t.Name,
	}

	slog.Info(t.Code(), logAttrs...)

	p.memory.Update(t.ID, t.SubID, t.Status)
	p.memory.wgTasks.Done()
}
