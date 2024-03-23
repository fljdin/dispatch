package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Process struct {
	ID      int
	memory  *Memory
	context context.Context
}

func NewProcess(memory *Memory, ctx context.Context) *Process {
	return &Process{
		memory:  memory,
		context: ctx,
	}
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
			ProcID:  p.ID,
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

	p.memory.results <- tasks.Result{
		ID:        t.ID,
		SubID:     t.SubID,
		Name:      t.Name,
		Action:    t.Action.String(),
		ProcID:    p.ID,
		Status:    report.Status,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
		Elapsed:   report.Elapsed,
		Output:    report.Output,
		Error:     report.Error,
	}
}
