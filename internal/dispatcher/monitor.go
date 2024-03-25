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
			m.memory.SetStatus(result.ID, result.SubID, result.Status)
			m.memory.wgTasks.Done()
		}
	}
}
