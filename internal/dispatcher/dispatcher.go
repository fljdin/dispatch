package dispatcher

import (
	"context"
	"os"

	"github.com/fljdin/dispatch/internal/models"
)

type Dispatcher struct {
	cancel  func()
	context context.Context
	memory  *SharedMemory
}

func NewDispatcher(ctx context.Context, count int, size int) *Dispatcher {
	ctx, cancel := context.WithCancel(ctx)

	d := &Dispatcher{
		context: ctx,
		cancel:  cancel,
	}

	d.memory = &SharedMemory{
		tasks:   make(chan models.Task, size),
		results: make(chan models.TaskResult, size),
	}

	d.launchOberserver()
	d.launchWorkers(count, size)

	return d
}

func (d *Dispatcher) launchOberserver() {
	observer := &Observer{
		memory:  d.memory,
		context: d.context,
	}

	go observer.Start()
}

func (d *Dispatcher) launchWorkers(count int, size int) {
	for i := 1; i <= count; i++ {
		worker := &Worker{
			ID:      i,
			memory:  d.memory,
			context: d.context,
		}
		go worker.Start()
	}
}

func (d *Dispatcher) AddTask(task models.Task) {
	d.memory.AddTask(task)
}

func (d *Dispatcher) GetStatus(ID int) int {
	return d.memory.GetStatus(ID)
}

func (d *Dispatcher) TraceTo(filename string) error {
	var err error
	const flag int = os.O_APPEND | os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	if len(filename) > 0 {
		d.memory.trace, err = os.OpenFile(filename, flag, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Dispatcher) Wait() {
	d.memory.wgTasks.Wait()   // wait until each task has been processed
	d.cancel()                // warm workers to stop theirs loop
	d.memory.wgWorkers.Wait() // wait until each worker has been stopped
}
