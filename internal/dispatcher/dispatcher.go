package dispatcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fljdin/dispatch/internal/models"
)

type Dispatcher struct {
	statuses   models.StatusMap
	context    context.Context
	tasks      chan models.Task
	results    chan models.TaskResult
	wgTasks    sync.WaitGroup
	wgWorkers  sync.WaitGroup
	wgObserver sync.WaitGroup
	cancel     func()
	trace      *os.File
}

func NewDispatcher(ctx context.Context, count int, size int) *Dispatcher {
	ctx, cancel := context.WithCancel(ctx)

	d := &Dispatcher{
		context:  ctx,
		cancel:   cancel,
		tasks:    make(chan models.Task, size),
		results:  make(chan models.TaskResult, size),
		statuses: models.StatusMap{},
	}

	launchObserver(d)
	launchWorkers(count, d)

	return d
}

func launchWorkers(count int, d *Dispatcher) {
	for i := 1; i <= count; i++ {
		worker := &Worker{
			ID:         i,
			dispatcher: d,
		}
		go worker.Start()
		d.wgWorkers.Add(1)
	}
}

func launchObserver(d *Dispatcher) {
	go d.observer(d.context)
	d.wgObserver.Add(1)
}

func (d *Dispatcher) Add(task models.Task) {
	d.tasks <- task
	d.wgTasks.Add(1)
}

func (d *Dispatcher) TraceTo(filename string) error {
	var err error
	const flag int = os.O_APPEND | os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	if len(filename) > 0 {
		d.trace, err = os.OpenFile(filename, flag, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Dispatcher) GetStatus(ID int) int {
	return d.statuses.Load(ID)
}

func (d *Dispatcher) Wait() {
	d.wgTasks.Wait()    // wait until each task has been processed
	d.cancel()          // warm workers to stop theirs loop
	d.wgWorkers.Wait()  // wait until each worker has been stopped
	d.wgObserver.Wait() // wait until observer has been stopped
}

func (d *Dispatcher) observer(ctx context.Context) {
	defer d.wgObserver.Done()
	defer d.trace.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-d.results:
			d.statuses.Store(result.ID, result.Status)
			d.logger(result)
			d.tracer(result)
			d.wgTasks.Done()
		}
	}
}

func (d *Dispatcher) logger(result models.TaskResult) {
	log.Printf(
		"Worker %d completed Task %d (query #%d) (success: %t, elapsed: %s)\n",
		result.WorkerID,
		result.ID,
		result.QueryID,
		(result.Status == models.Succeeded),
		result.Elapsed.Round(time.Millisecond),
	)
}

func (d *Dispatcher) tracer(result models.TaskResult) {
	if d.trace != nil {
		template := `===== Task %d (query #%d) (success: %t, elapsed: %s) =====
Started at: %s
Ended at:   %s
Error: %s
Output:
%s
`
		report := fmt.Sprintf(
			template,
			result.ID,
			result.QueryID,
			(result.Status == models.Succeeded),
			result.Elapsed.Round(time.Millisecond),
			result.StartTime.String(),
			result.EndTime.String(),
			result.Error,
			result.Output,
		)

		_, err := d.trace.Write([]byte(report))
		if err != nil {
			log.Println(err)
		}
	}
}
