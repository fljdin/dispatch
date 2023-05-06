package dispatcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fljdin/dispatch/internal/models"
)

type Observer struct {
	mem *WorkerMem
	ctx context.Context
}

func (o *Observer) Start() {
	o.mem.wgWorkers.Add(1)
	defer o.mem.wgWorkers.Done()
	defer o.mem.trace.Close()

	for {
		select {
		case <-o.ctx.Done():
			return
		case result := <-o.mem.results:
			o.mem.statuses.Store(result.ID, result.Status)
			o.logger(result)
			o.tracer(result)
			o.mem.wgTasks.Done()
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
	if o.mem.trace != nil {
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

		_, err := o.mem.trace.Write([]byte(report))
		if err != nil {
			log.Println(err)
		}
	}
}
