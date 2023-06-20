package dispatcher

import (
	"context"
	"sync"

	"github.com/fljdin/dispatch/internal/task"
)

type Dispatcher struct {
	cancel   func()
	context  context.Context
	workers  int
	observer *Observer
	memory   *Memory
}

func (d *Dispatcher) Wait() {
	d.launchObserver()
	d.launchWorkers()

	d.memory.wgTasks.Wait()
	d.cancel()
	d.memory.wgWorkers.Wait()
}

func (d *Dispatcher) AddTask(task task.Task) {
	d.memory.AddTask(task)
}

func (d *Dispatcher) GetStatus(id int) int {
	return d.memory.GetStatus(id)
}

func (d *Dispatcher) launchObserver() {
	go d.observer.Start()
}

func (d *Dispatcher) launchWorkers() {
	for i := 1; i <= d.workers; i++ {
		worker := &Worker{
			ID:      i,
			memory:  d.memory,
			context: d.context,
		}
		go worker.Start()
	}
}

type Memory struct {
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	queue     task.TaskQueue
	results   chan task.TaskResult
}

func (m *Memory) GetStatus(id int) int {
	return m.queue.GetStatus(id)
}

func (m *Memory) SetStatus(id int, status int) {
	m.queue.SetStatus(id, status)
}

func (m *Memory) AddTask(task task.Task) {
	m.queue.Push(&task)
	m.wgTasks.Add(1)
}

func (m *Memory) ForwardTask(task task.Task) {
	m.queue.Push(&task)
}

func (m *Memory) StartWorker() {
	m.wgWorkers.Add(1)
}

func (m *Memory) EndWorker() {
	m.wgWorkers.Done()
}
