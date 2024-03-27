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
	wgTasks sync.WaitGroup
	wgProcs sync.WaitGroup
	queue   queue.Queue
	tasks   chan tasks.Task
	results chan Result
}

func (m *Memory) Evaluate(id int) status.Status {
	return m.queue.Evaluate(id)
}

func (m *Memory) AddTask(task tasks.Task) {
	m.queue.Add(task)
	m.wgTasks.Add(1)
}

func (m *Memory) Done(tid tasks.TaskIdentifier, status status.Status) {
	m.queue.Update(tid, status)
	m.wgTasks.Done()
}

func (m *Memory) SendTask(task tasks.Task) {
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
