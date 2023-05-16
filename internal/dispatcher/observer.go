package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/logger"
	"github.com/fljdin/dispatch/internal/models"
)

type Observer struct {
	memory  *SharedMemory
	context context.Context
	loggers []logger.Logger
}

func (o *Observer) Start() {
	o.memory.StartWorker()
	defer o.memory.EndWorker()

	for {
		select {
		case <-o.context.Done():
			return
		case result := <-o.memory.results:
			o.memory.SetStatus(result.ID, result.Status)
			o.Log(result)
			o.memory.wgTasks.Done()
		}
	}
}

func (o *Observer) WithConsole() {
	o.loggers = append(o.loggers, &logger.Console{})
}

func (o *Observer) WithTrace(filename string) error {

	if len(filename) > 0 {
		trace := &logger.Trace{
			Filename: filename,
		}

		if err := trace.Open(); err != nil {
			return err
		}

		o.loggers = append(o.loggers, trace)
	}

	return nil
}

func (o *Observer) Log(result models.TaskResult) {
	for _, logger := range o.loggers {
		logger.Render(result)
	}
}
