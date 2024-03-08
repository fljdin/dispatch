package tasks

import (
	"container/list"
	"sync"

	"golang.org/x/exp/slices"
)

type Queue struct {
	mut    sync.Mutex
	status *StatusMap
	tasks  *list.List
}

func NewQueue() Queue {
	return Queue{
		status: NewStatusMap(),
		tasks:  list.New(),
	}
}

func (q *Queue) Status(taskId int) int {
	return q.status.Get(taskId)
}

func (q *Queue) SetStatus(taskId, taskSubId, status int) {
	q.status.Set(taskId, taskSubId, status)
}

func (q *Queue) Len() int {
	q.mut.Lock()
	defer q.mut.Unlock()

	return q.tasks.Len()
}

func (q *Queue) Add(t Task) {
	q.mut.Lock()
	defer q.mut.Unlock()

	q.tasks.PushBack(t)
	q.status.Set(t.ID, t.SubID, t.Status)
}

func (q *Queue) Pop() (Task, bool) {
	q.mut.Lock()
	defer q.mut.Unlock()

	element := q.tasks.Front()
	if element == nil {
		return Task{}, false
	}

	task := element.Value.(Task)
	q.tasks.Remove(element)

	task.Status = q.evaluate(task)
	return task, true
}

func (q *Queue) evaluate(t Task) int {
	for _, id := range t.Depends {
		status := q.status.Get(id)

		if slices.Contains([]int{Interrupted, Failed}, status) {
			return Interrupted
		} else if status != Succeeded {
			return Waiting
		}
	}

	return Ready
}
