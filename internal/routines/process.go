package routines

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/status"
)

type Process struct {
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
			p.run(task)
		}
	}
}

func (p *Process) run(t config.Task) {
	if t.Status == status.Interrupted {
		slog.Error(
			fmt.Sprintf("task=%s", t),
			"status", status.Interrupted,
			"name", t.Name,
		)

		p.memory.Done(t.Identifier, t.Status)
		return
	}

	report, commands := t.Action.Run()
	for id, command := range commands {
		p.memory.AddTask(config.Task{
			Identifier: config.NewId(t.Identifier.ID, id+1),
			Action:     command,
			Name:       t.Name,
		})
	}

	var slogFunc func(string, ...any) = slog.Info
	if report.Status.IsFailed() {
		slogFunc = slog.Error
	}

	slogFunc(
		fmt.Sprintf("task=%s", t),
		"status", report.Status,
		"name", t.Name,
		"elapsed", report.Elapsed.Round(time.Millisecond),
	)

	slogFunc(fmt.Sprintf("task=%s cmd=%s action: %s", t, t.Action.Command(), t.Action.String()))

	slog.Debug(
		fmt.Sprintf("task=%s", t),
		"start", report.StartTime.Format(time.DateTime),
		"end", report.EndTime.Format(time.DateTime),
	)

	if len(report.Error) > 0 {
		slogFunc(fmt.Sprintf("task=%s error: %s", t, report.Error))
	}

	if len(report.Output) > 0 {
		slogFunc(fmt.Sprintf("task=%s output: %s", t, report.Output))
	}

	p.memory.Done(t.Identifier, report.Status)
}
