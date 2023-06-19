package models

import (
	"container/list"
)

type TaskQueue struct {
	tasks    *list.List
	statuses StatusMap
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks: list.New(),
	}
}

func (q *TaskQueue) GetStatus(ID int) int {
	return q.statuses.Load(ID)
}

func (q *TaskQueue) SetStatus(ID int, status int) {
	q.statuses.Store(ID, status)
}

func (q *TaskQueue) Len() int {
	return q.tasks.Len()
}

func (q *TaskQueue) Push(task *Task) {
	q.tasks.PushBack(task)
	q.statuses.Store(task.ID, task.Status)
}

func (q *TaskQueue) Pop() *Task {
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
		parentStatus := q.statuses.Load(id)

		if parentStatus >= Failed {
			return Interrupted
		} else if parentStatus < Succeeded {
			return Waiting
		}
	}

	return Ready
}
