package dispatcher

import (
	"sync"

	"github.com/fljdin/dispatch/internal/models"
)

type SharedMemory struct {
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	statuses  models.StatusMap
	tasks     chan models.Task
	results   chan models.TaskResult
}

func (s *SharedMemory) AddTask(task models.Task) {
	s.tasks <- task
	s.wgTasks.Add(1)
}

func (s *SharedMemory) ForwardTask(task models.Task) {
	s.tasks <- task
}

func (s *SharedMemory) StartWorker() {
	s.wgWorkers.Add(1)
}

func (s *SharedMemory) EndWorker() {
	s.wgWorkers.Done()
}

func (s *SharedMemory) GetStatus(ID int) int {
	return s.statuses.Load(ID)
}

func (s *SharedMemory) SetStatus(ID int, status int) {
	s.statuses.Store(ID, status)
}
