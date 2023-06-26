package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	. "github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/require"
)

func TestQueuePush(t *testing.T) {
	r := require.New(t)

	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "echo test",
		},
	})

	r.Equal(1, queue.Len())
}

func TestQueuePopEmpty(t *testing.T) {
	r := require.New(t)

	queue := NewQueue()
	_, ok := queue.Pop()

	r.Equal(false, ok)
}

func TestQueuePop(t *testing.T) {
	r := require.New(t)

	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "echo test",
		},
	})
	task, _ := queue.Pop()

	r.Equal(1, task.ID)
	r.Equal(0, queue.Len())
}

func TestQueueAutonomousTaskMustBeReady(t *testing.T) {
	r := require.New(t)

	queue := NewQueue()
	queue.Add(Task{
		ID: 1,
		Action: Command{
			Text: "echo test",
		},
	})
	task, _ := queue.Pop()

	r.Equal(Ready, task.Status)
}

func TestQueueDependentTaskMustBeWaiting(t *testing.T) {
	r := require.New(t)

	queue := NewQueue()
	queue.Add(Task{
		ID:      2,
		Depends: []int{1},
		Action: Command{
			Text: "true",
		},
	})
	task, _ := queue.Pop()

	r.Equal(Waiting, task.Status)
}

func TestQueueDependentTaskMustBeReady(t *testing.T) {
	r := require.New(t)

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

	r.Equal(Ready, task.Status)
}

func TestQueueDependentTaskMustBeInterrupted(t *testing.T) {
	r := require.New(t)

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

	r.Equal(Interrupted, task.Status)
}
