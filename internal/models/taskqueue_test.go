package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueuePush(t *testing.T) {
	queue := NewTaskQueue()
	queue.Push(&Task{
		ID: 1,
		Command: Command{
			Text: "echo test",
		},
	})

	assert.Equal(t, 1, queue.Len())
}

func TestQueuePop(t *testing.T) {
	queue := NewTaskQueue()
	queue.Push(&Task{
		ID: 1,
		Command: Command{
			Text: "echo test",
		},
	})

	task := queue.Pop()
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, 0, queue.Len())
}
