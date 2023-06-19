package models

import "container/list"

type TaskQueue struct {
	tasks *list.List
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks: list.New(),
	}
}

func (q *TaskQueue) Len() int {
	return q.tasks.Len()
}

func (q *TaskQueue) Push(task *Task) {
	q.tasks.PushBack(task)
}

func (q *TaskQueue) Pop() *Task {
	element := q.tasks.Front()
	if element == nil {
		return nil
	}

	task := element.Value.(*Task)
	q.tasks.Remove(element)
	return task
}
