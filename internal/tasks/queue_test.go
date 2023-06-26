package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	. "github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/assert"
)

func TestQueuePush(t *testing.T) {
	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "echo test",
		},
	})

	assert.Equal(t, 1, queue.Len())
}

func TestQueuePopEmpty(t *testing.T) {
	queue := NewQueue()
	_, ok := queue.Pop()
	assert.Equal(t, false, ok)
}

func TestQueuePop(t *testing.T) {
	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "echo test",
		},
	})

	task, _ := queue.Pop()
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, 0, queue.Len())
}

func TestQueueAutonomousTaskMustBeReady(t *testing.T) {
	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "echo test",
		},
	})

	task, _ := queue.Pop()
	assert.Equal(t, Ready, task.Status)
}

func TestQueueDependentTaskMustBeWaiting(t *testing.T) {
	queue := NewQueue()
	queue.Add(Task{
		ID:      2,
		Depends: []int{1},
		Action: Command{
			Text: "true",
		},
	})

	task, _ := queue.Pop()
	assert.Equal(t, Waiting, task.Status)
}

func TestQueueDependentTaskMustBeReady(t *testing.T) {
	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "true",
		},
	})
	queue.Add(Task{
		ID:      2,
		Depends: []int{1},
		Action: Command{
			Text: "true",
		},
	})

	_, _ = queue.Pop()
	queue.SetStatus(1, Succeeded)

	task, _ := queue.Pop()
	assert.Equal(t, Ready, task.Status)
}

func TestQueueDependentTaskMustBeInterrupted(t *testing.T) {
	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "true",
		},
	})
	queue.Add(Task{
		ID:      2,
		Depends: []int{1},
		Action: Command{
			Text: "true",
		},
	})

	_, _ = queue.Pop()
	queue.SetStatus(1, Interrupted)

	task, _ := queue.Pop()
	assert.Equal(t, Interrupted, task.Status)
}
