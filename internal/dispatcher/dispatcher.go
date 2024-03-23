package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Dispatcher struct {
	cancel  func()
	context context.Context
	workers int
	monitor *Monitor
	memory  *Memory
}

func (d *Dispatcher) Wait() {
	d.launchMonitor()
	d.launchWorkers()

	d.memory.WaitForTasks()
	d.cancel()
	d.memory.WaitForWorkers()
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

func (d *Dispatcher) launchWorkers() {
	for i := 1; i <= d.workers; i++ {
		worker := &Worker{
			ID:      i,
			memory:  d.memory,
			context: d.context,
		}
		go worker.Start()
	}
}
