package tasks

import (
	"container/list"
	"sync"
)

type Queue struct {
	tasks  *list.List
	status StatusMap
	mut    sync.Mutex
}

func NewQueue() Queue {
	return Queue{
		tasks: list.New(),
	}
}

func (q *Queue) GetStatus(id int) int {
	return q.status.Get(id)
}

func (q *Queue) SetStatus(id int, status int) {
	q.status.Set(id, status)
}

func (q *Queue) Len() int {
	q.mut.Lock()
	defer q.mut.Unlock()

	return q.tasks.Len()
}

func (q *Queue) Push(task *Task) {
	q.mut.Lock()
	defer q.mut.Unlock()

	q.tasks.PushBack(task)
	q.status.Set(task.ID, task.Status)
}

func (q *Queue) Pop() *Task {
	q.mut.Lock()
	defer q.mut.Unlock()

	element := q.tasks.Front()
	if element == nil {
		return nil
	}

	task := element.Value.(*Task)
	q.tasks.Remove(element)

	task.Status = q.evaluate(task)
	return task
}

func (q *Queue) evaluate(task *Task) int {
	for _, id := range task.Depends {
		parentStatus := q.status.Get(id)

		if parentStatus >= Failed {
			return Interrupted
		} else if parentStatus < Succeeded {
			return Waiting
		}
	}

	return Ready
}
