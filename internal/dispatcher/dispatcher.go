package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Dispatcher struct {
	cancel    func()
	context   context.Context
	processes int
	monitor   *Monitor
	memory    *Memory
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
