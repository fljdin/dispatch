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
	msg := fmt.Sprintf("task=%s", t)

	for id, command := range commands {
		p.memory.AddTask(tasks.Task{
			Identifier: tasks.NewId(t.Identifier.ID, id+1),
			Action:     command,
			Name:       t.Name,
		})
	}

	logAttrs := []any{
		"status", report.Status,
		"name", t.Name,
		"start", report.StartTime.Format(time.DateTime),
		"elapsed", report.Elapsed.Round(time.Millisecond),
	}

	if report.Status.IsFailed() {
		slog.Error(msg, logAttrs...)
	} else {
		slog.Info(msg, logAttrs...)
	}

	slog.Info(fmt.Sprintf("%s action:\n%s", msg, t.Action.String()))

	if len(report.Error) > 0 {
		slog.Error(fmt.Sprintf("%s error:\n%s", msg, report.Error))
	}

	if len(report.Output) > 0 {
		slog.Info(fmt.Sprintf("%s output:\n%s", msg, report.Output))
	}

	p.memory.results <- Result{
		Identifier: t.Identifier,
		Status:     report.Status,
	}
}

func (p *Process) interrupt(t tasks.Task) {
	logAttrs := []any{
		"status", status.Interrupted,
		"name", t.Name,
	}

	slog.Info(t.String(), logAttrs...)

	p.memory.Update(t.Identifier, t.Status)
	p.memory.wgTasks.Done()
}
