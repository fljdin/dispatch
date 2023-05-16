package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/models"
)

type Dispatcher struct {
	cancel   func()
	context  context.Context
	workers  int
	observer *Observer
	memory   *SharedMemory
}

func (d *Dispatcher) Wait() {
	d.launchObserver()
	d.launchWorkers()

	d.memory.wgTasks.Wait()
	d.cancel()
	d.memory.wgWorkers.Wait()
}

func (d *Dispatcher) AddTask(task models.Task) {
	d.memory.AddTask(task)
}

func (d *Dispatcher) GetStatus(ID int) int {
	return d.memory.GetStatus(ID)
}

func (d *Dispatcher) launchObserver() {
	go d.observer.Start()
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
