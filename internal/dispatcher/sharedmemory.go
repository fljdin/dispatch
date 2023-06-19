package dispatcher

import (
	"sync"

	"github.com/fljdin/dispatch/internal/models"
)

type SharedMemory struct {
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	queue     models.TaskQueue
	results   chan models.TaskResult
}

func (s *SharedMemory) AddTask(task models.Task) {
	s.queue.Push(&task)
	s.wgTasks.Add(1)
}

func (s *SharedMemory) ForwardTask(task models.Task) {
	s.queue.Push(&task)
}

func (s *SharedMemory) StartWorker() {
	s.wgWorkers.Add(1)
}

func (s *SharedMemory) EndWorker() {
	s.wgWorkers.Done()
}

func (s *SharedMemory) GetStatus(ID int) int {
	return s.queue.GetStatus(ID)
}

func (s *SharedMemory) SetStatus(ID int, status int) {
	s.queue.SetStatus(ID, status)
}
