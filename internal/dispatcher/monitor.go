package dispatcher

import (
	"context"
	"fmt"
	"log/slog"
)

type Monitor struct {
	memory  *Memory
	context context.Context
}

func (m *Monitor) Start() {
	m.memory.StartProcess()
	defer m.memory.EndProcess()

	for {
		select {
		case <-m.context.Done():
			return
		case result := <-m.memory.results:
			m.memory.Update(result.Identifier, result.Status)
			m.memory.wgTasks.Done()

			// fill back the tasks channel
			if task, ok := m.memory.queue.Next(); ok {
				m.memory.tasks <- task
				slog.Debug(
					fmt.Sprintf("task=%s", task),
					"msg", "task sent to internal channel",
				)
			}
		}
	}
}
