package dispatcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fljdin/dispatch/internal/models"
)

type Observer struct {
	memory  *SharedMemory
	context context.Context
	trace   *os.File
}

func (o *Observer) Start() {
	o.memory.StartWorker()
	defer o.memory.EndWorker()
	defer o.trace.Close()

	for {
		select {
		case <-o.context.Done():
			return
		case result := <-o.memory.results:
			o.memory.statuses.Store(result.ID, result.Status)
			o.Log(result)
			o.memory.wgTasks.Done()
		}
	}
}

func (o *Observer) TraceTo(filename string) error {
	var err error
	const flag int = os.O_APPEND | os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	if len(filename) > 0 {
		o.trace, err = os.OpenFile(filename, flag, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Observer) Log(result models.TaskResult) {
	o.console(result)
	o.tracer(result)
}

func (o *Observer) console(result models.TaskResult) {
	log.Printf(
		"Worker %d completed Task %d (query #%d) (success: %t, elapsed: %s)\n",
		result.WorkerID,
		result.ID,
		result.QueryID,
		(result.Status == models.Succeeded),
		result.Elapsed.Round(time.Millisecond),
	)
}

func (o *Observer) tracer(result models.TaskResult) {
	if o.trace != nil {
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

		_, err := o.trace.Write([]byte(report))
		if err != nil {
			log.Println(err)
		}
	}
}
