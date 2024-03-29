package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Dispatcher struct {
	cancel  func()
	context context.Context
	*Memory
}

func New(procs int) Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	d := Dispatcher{
		context: ctx,
		cancel:  cancel,
		Memory: &Memory{
			processes: procs,
			queue:     NewQueue(),
			tasks:     make(chan tasks.Task, procs),
		},
	}

	return d
}

func (d *Dispatcher) Wait() {
	d.FillTasks()
	d.launchProcesses()
	d.WaitForTasks()
	d.cancel()
	d.WaitForProcesses()
}

func (d Dispatcher) launchProcesses() {
	for i := 1; i <= d.processes; i++ {
		process := Process{
			memory:  d.Memory,
			context: d.context,
		}
		go process.Start()
	}
}
