package task

import (
	"container/list"
	"sync"
)

type TaskQueue struct {
	tasks  *list.List
	status StatusMap
	mut    sync.Mutex
}

func NewQueue() TaskQueue {
	return TaskQueue{
		tasks: list.New(),
	}
}

func (q *TaskQueue) GetStatus(id int) int {
	return q.status.Get(id)
}

func (q *TaskQueue) SetStatus(id int, status int) {
	q.status.Set(id, status)
}

func (q *TaskQueue) Len() int {
	q.mut.Lock()
	defer q.mut.Unlock()

	return q.tasks.Len()
}

func (q *TaskQueue) Push(task *Task) {
	q.mut.Lock()
	defer q.mut.Unlock()

	q.tasks.PushBack(task)
	q.status.Set(task.ID, task.Status)
}

func (q *TaskQueue) Pop() *Task {
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

func (q *TaskQueue) evaluate(task *Task) int {
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
