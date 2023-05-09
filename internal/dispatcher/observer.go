package dispatcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fljdin/dispatch/internal/models"
)

type Observer struct {
	memory  *SharedMemory
	context context.Context
}

func (o *Observer) Start() {
	o.memory.StartWorker()
	defer o.memory.EndWorker()
	defer o.memory.trace.Close()

	for {
		select {
		case <-o.context.Done():
			return
		case result := <-o.memory.results:
			o.memory.statuses.Store(result.ID, result.Status)
			o.logger(result)
			o.tracer(result)
			o.memory.wgTasks.Done()
		}
	}
}

func (o *Observer) logger(result models.TaskResult) {
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
	if o.memory.trace != nil {
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

		_, err := o.memory.trace.Write([]byte(report))
		if err != nil {
			log.Println(err)
		}
	}
}
