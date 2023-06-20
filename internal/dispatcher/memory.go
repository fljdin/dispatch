package dispatcher

import (
	"sync"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Memory struct {
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	queue     tasks.Queue
	results   chan tasks.Result
}

func (m *Memory) GetStatus(id int) int {
	return m.queue.GetStatus(id)
}

func (m *Memory) SetStatus(id int, status int) {
	m.queue.SetStatus(id, status)
}

func (m *Memory) AddTask(task tasks.Task) {
	m.queue.Add(&task)
	m.wgTasks.Add(1)
}

func (m *Memory) ForwardTask(task tasks.Task) {
	m.queue.Add(&task)
}

func (m *Memory) StartWorker() {
	m.wgWorkers.Add(1)
}

func (m *Memory) EndWorker() {
	m.wgWorkers.Done()
}
