package dispatcher

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/dispatch/internal/tasks"
)

type Result struct {
	Identifier tasks.TaskIdentifier
	Status     status.Status
}

type Memory struct {
	active    int
	processes int
	queue     queue.Queue
	results   chan Result
	tasks     chan tasks.Task
	wgProcs   sync.WaitGroup
	wgTasks   sync.WaitGroup
}

func (m *Memory) Evaluate(id int) status.Status {
	return m.queue.Evaluate(id)
}

func (m *Memory) AddTask(task tasks.Task) {
	m.queue.Add(task)
	m.wgTasks.Add(1)
}

func (m *Memory) Done(tid tasks.TaskIdentifier, status status.Status) {
	m.active--
	m.queue.Update(tid, status)
	m.wgTasks.Done()
}

// fill back the tasks channel for any idle processes
func (m *Memory) FillTasks() {
	idleProcs := m.processes - m.active
	slog.Debug("filling tasks channel", "idle", idleProcs)

	for i := 0; i < idleProcs; i++ {
		if task, ok := m.queue.Next(); ok {
			m.queuing(task)
		}
	}
}

func (m *Memory) queuing(task tasks.Task) {
	m.active++
	m.queue.Update(task.Identifier, status.Running)
	slog.Debug(
		fmt.Sprintf("task=%s", task),
		"msg", "task sent to internal channel",
	)

	m.tasks <- task
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
