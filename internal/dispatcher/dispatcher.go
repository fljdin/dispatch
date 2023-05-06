package dispatcher

import (
	"context"
	"os"
	"sync"

	"github.com/fljdin/dispatch/internal/models"
)

type WorkerMem struct {
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	statuses  models.StatusMap
	tasks     chan models.Task
	results   chan models.TaskResult
	trace     *os.File
}

type Dispatcher struct {
	cancel  func()
	context context.Context
	workMem *WorkerMem
}

func NewDispatcher(ctx context.Context, count int, size int) *Dispatcher {
	ctx, cancel := context.WithCancel(ctx)

	d := &Dispatcher{
		context: ctx,
		cancel:  cancel,
	}

	d.launchWorkers(count, size)

	return d
}

func (d *Dispatcher) launchWorkers(count int, size int) {

	d.workMem = &WorkerMem{
		tasks:   make(chan models.Task, size),
		results: make(chan models.TaskResult, size),
	}

	observer := &Observer{
		mem: d.workMem,
		ctx: d.context,
	}

	go observer.Start()

	for i := 1; i <= count; i++ {
		worker := &Worker{
			ID:  i,
			mem: d.workMem,
			ctx: d.context,
		}
		go worker.Start()
	}
}

func (d *Dispatcher) Add(task models.Task) {
	d.workMem.tasks <- task
	d.workMem.wgTasks.Add(1)
}

func (d *Dispatcher) GetStatus(ID int) int {
	return d.workMem.statuses.Load(ID)
}

func (d *Dispatcher) TraceTo(filename string) error {
	var err error
	const flag int = os.O_APPEND | os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	if len(filename) > 0 {
		d.workMem.trace, err = os.OpenFile(filename, flag, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Dispatcher) Wait() {
	d.workMem.wgTasks.Wait()   // wait until each task has been processed
	d.cancel()                 // warm workers to stop theirs loop
	d.workMem.wgWorkers.Wait() // wait until each worker has been stopped
}
