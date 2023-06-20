package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/assert"
)

func TestQueuePush(t *testing.T) {
	queue := NewQueue()
	queue.Push(&Task{
		ID: 1,
		Command: Command{
			Text: "echo test",
		},
	})

	assert.Equal(t, 1, queue.Len())
}

func TestQueuePopEmpty(t *testing.T) {
	queue := NewQueue()
	task := queue.Pop()
	assert.Nil(t, task)
}

func TestQueuePop(t *testing.T) {
	queue := NewQueue()
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

func TestQueueAutonomousTaskMustBeReady(t *testing.T) {
	queue := NewQueue()
	queue.Push(&Task{
		ID: 1,
		Command: Command{
			Text: "echo test",
		},
	})

	task := queue.Pop()
	assert.Equal(t, Ready, task.Status)
}

func TestQueueDependentTaskMustBeWaiting(t *testing.T) {
	queue := NewQueue()
	queue.Push(&Task{
		ID:      2,
		Depends: []int{1},
		Command: Command{
			Text: "true",
		},
	})

	task := queue.Pop()
	assert.Equal(t, Waiting, task.Status)
}

func TestQueueDependentTaskMustBeReady(t *testing.T) {
	queue := NewQueue()
	queue.Push(&Task{
		ID: 1,
		Command: Command{
			Text: "true",
		},
	})
	queue.Push(&Task{
		ID:      2,
		Depends: []int{1},
		Command: Command{
			Text: "true",
		},
	})

	_ = queue.Pop()
	queue.SetStatus(1, Succeeded)

	task := queue.Pop()
	assert.Equal(t, Ready, task.Status)
}

func TestQueueDependentTaskMustBeInterrupted(t *testing.T) {
	queue := NewQueue()
	queue.Push(&Task{
		ID: 1,
		Command: Command{
			Text: "true",
		},
	})
	queue.Push(&Task{
		ID:      2,
		Depends: []int{1},
		Command: Command{
			Text: "true",
		},
	})

	_ = queue.Pop()
	queue.SetStatus(1, Interrupted)

	task := queue.Pop()
	assert.Equal(t, Interrupted, task.Status)
}
