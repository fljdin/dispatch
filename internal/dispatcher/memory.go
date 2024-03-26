package dispatcher

import (
	"sync"

	"github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/dispatch/internal/tasks"
)

type Memory struct {
	wgTasks sync.WaitGroup
	wgProcs sync.WaitGroup
	queue   queue.Queue
	tasks   chan tasks.Task
	results chan Result
}

func (m *Memory) Evaluate(id int) status.Status {
	return m.queue.Evaluate(id)
}

func (m *Memory) Update(id, subid int, status status.Status) {
	m.queue.Update(id, subid, status)
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
