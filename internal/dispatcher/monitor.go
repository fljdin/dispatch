package dispatcher

import (
	"context"

	"github.com/fljdin/dispatch/internal/logger"
	"github.com/fljdin/dispatch/internal/tasks"
)

type Monitor struct {
	memory  *Memory
	context context.Context
	loggers []logger.Logger
}

func NewMonitor(memory *Memory, ctx context.Context) *Monitor {
	return &Monitor{
		memory:  memory,
		context: ctx,
	}
}

func (m *Monitor) Start() {
	m.memory.StartWorker()
	defer m.memory.EndWorker()

	for {
		select {
		case <-m.context.Done():
			return
		case result := <-m.memory.results:
			m.memory.SetStatus(result.ID, result.SubID, result.Status)
			m.Log(result)
			m.memory.wgTasks.Done()
		}
	}
}

func (m *Monitor) WithConsole() {
	m.loggers = append(m.loggers, &logger.Console{})
}

func (m *Monitor) WithTrace(filename string) error {

	if len(filename) > 0 {
		trace := &logger.Trace{
			Filename: filename,
		}

		if err := trace.Open(); err != nil {
			return err
		}

		m.loggers = append(m.loggers, trace)
	}

	return nil
}

func (m Monitor) Log(result tasks.Result) {
	for _, logger := range m.loggers {
		logger.Render(result)
	}
}
