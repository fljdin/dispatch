package dispatcher

import (
	"sync"

	"github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/tasks"
)

type Memory struct {
	wgTasks sync.WaitGroup
	wgProcs sync.WaitGroup
	queue   queue.Queue
	results chan Result
}

func NewMemory() *Memory {
	return &Memory{
		queue:   queue.New(),
		results: make(chan Result),
	}
}

func (m *Memory) Status(taskId int) int {
	return m.queue.Status(taskId)
}

func (m *Memory) SetStatus(taskId, taskSubId, status int) {
	m.queue.SetStatus(taskId, taskSubId, status)
}

func (m *Memory) AddTask(task tasks.Task) {
	m.queue.Add(task)
	m.wgTasks.Add(1)
}

func (m *Memory) ForwardTask(task tasks.Task) {
	m.queue.Add(task)
}

func (m *Memory) StartProcess() {
	m.wgProcs.Add(1)
}

func (m *Memory) EndProcess() {
	m.wgProcs.Done()
}

func (m *Memory) WaitForTasks() {
	m.wgTasks.Wait()
}

func (m *Memory) WaitForProcesses() {
	m.wgProcs.Wait()
}
