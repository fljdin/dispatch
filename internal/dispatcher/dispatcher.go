package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/tasks"
)

type Dispatcher struct {
	cancel    func()
	context   context.Context
	processes int
	monitor   *Monitor
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
			results: make(chan tasks.Result, 10),
		},
	}

	d.monitor = NewMonitor(
		d.memory,
		d.context,
	)

	return d
}

func (d *Dispatcher) Wait() {
	d.launchMonitor()
	d.launchProcesses()

	d.memory.WaitForTasks()
	d.cancel()
	d.memory.WaitForProcesses()
}

func (d *Dispatcher) AddTask(task tasks.Task) {
	d.memory.AddTask(task)
}

func (d *Dispatcher) Status(taskId int) int {
	return d.memory.Status(taskId)
}

func (d *Dispatcher) launchMonitor() {
	go d.monitor.Start()
}

func (d *Dispatcher) launchProcesses() {
	for i := 1; i <= d.processes; i++ {
		process := &Process{
			ID:      i,
			memory:  d.memory,
			context: d.context,
		}
		go process.Start()
	}
}
