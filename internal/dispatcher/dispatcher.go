package dispatcher

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/dispatch/internal/tasks"
)

type Dispatcher struct {
	cancel    func()
	context   context.Context
	processes int
	memory    *Memory
}

func New(procs int) Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	d := Dispatcher{
		context:   ctx,
		cancel:    cancel,
		processes: procs,
		memory: &Memory{
			queue:   queue.New(),
			tasks:   make(chan tasks.Task, procs),
			results: make(chan Result, procs),
		},
	}

	return d
}

func (d *Dispatcher) Wait() {
	// fill tasks channel with ready tasks
	for i := 0; i < d.processes; i++ {
		if task, ok := d.memory.queue.Next(); ok {
			d.memory.tasks <- task
			d.memory.queue.Update(task.Identifier, status.Running)
			slog.Debug(
				fmt.Sprintf("task=%s", task),
				"msg", "task sent to internal channel",
			)
		}
	}

	d.launchMonitor()
	d.launchProcesses()

	d.memory.WaitForTasks()
	d.cancel()
	d.memory.WaitForProcesses()
}

func (d *Dispatcher) AddTask(task tasks.Task) {
	d.memory.AddTask(task)
}

func (d Dispatcher) Evaluate(id int) status.Status {
	return d.memory.Evaluate(id)
}

func (d Dispatcher) launchMonitor() {
	monitor := Monitor{
		memory:  d.memory,
		context: d.context,
	}
	go monitor.Start()
}

func (d Dispatcher) launchProcesses() {
	for i := 1; i <= d.processes; i++ {
		process := Process{
			ID:      i,
			memory:  d.memory,
			context: d.context,
		}
		go process.Start()
	}
}
