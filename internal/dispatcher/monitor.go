package dispatcher

import (
	"context"
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
			m.memory.Done(result.Identifier, result.Status)

			// fill back the tasks channel
			if task, ok := m.memory.queue.Next(); ok {
				m.memory.SendTask(task)
			}
		}
	}
}
