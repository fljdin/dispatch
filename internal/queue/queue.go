package queue

import (
	"container/list"
	"sync"

	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/dispatch/internal/tasks"
)

type Queue struct {
	mut    sync.Mutex
	status *StatusMap
	tasks  *list.List
}

func New() Queue {
	return Queue{
		status: NewStatusMap(),
		tasks:  list.New(),
	}
}

func (q *Queue) Status(taskId int) status.Status {
	return q.status.Get(taskId)
}

func (q *Queue) SetStatus(taskId, taskSubId int, status status.Status) {
	q.status.Set(taskId, taskSubId, status)
}

func (q *Queue) Len() int {
	q.mut.Lock()
	defer q.mut.Unlock()

	return q.tasks.Len()
}

func (q *Queue) Add(t tasks.Task) {
	q.mut.Lock()
	defer q.mut.Unlock()

	q.tasks.PushBack(t)
	q.status.Set(t.ID, t.SubID, t.Status)
}

func (q *Queue) Pop() (tasks.Task, bool) {
	q.mut.Lock()
	defer q.mut.Unlock()

	element := q.tasks.Front()
	if element == nil {
		return tasks.Task{}, false
	}

	task := element.Value.(tasks.Task)
	q.tasks.Remove(element)

	task.Status = q.evaluate(task)
	return task, true
}

func (q *Queue) evaluate(t tasks.Task) status.Status {
	for _, id := range t.Depends {
		s := q.status.Get(id)

		if s.IsFailed() {
			return status.Interrupted
		} else if s != status.Succeeded {
			return status.Waiting
		}
	}

	return status.Ready
}
