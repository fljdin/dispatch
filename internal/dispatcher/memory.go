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

func (m *Memory) Status(taskId int) int {
	return m.queue.Status(taskId)
}

func (m *Memory) SetStatus(taskId int, status int) {
	m.queue.SetStatus(taskId, status)
}

func (m *Memory) AddTask(task tasks.Task) {
	m.queue.Add(task)
	m.wgTasks.Add(1)
}

func (m *Memory) ForwardTask(task tasks.Task) {
	m.queue.Add(task)
}

func (m *Memory) StartWorker() {
	m.wgWorkers.Add(1)
}

func (m *Memory) EndWorker() {
	m.wgWorkers.Done()
}