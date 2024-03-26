package queue_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

func TestQueuePush(t *testing.T) {
	r := require.New(t)

	queue := New()
	queue.Add(tasks.Task{
		ID: 1,
		Action: tasks.Command{
			Text: "echo test",
		},
	})

	r.Equal(1, queue.Len())
}

func TestQueuePopEmpty(t *testing.T) {
	r := require.New(t)

	queue := New()
	_, ok := queue.Pop()

	r.Equal(false, ok)
}

func TestQueuePop(t *testing.T) {
	r := require.New(t)

	queue := New()
	queue.Add(tasks.Task{
		ID: 1,
		Action: tasks.Command{
			Text: "echo test",
		},
	})
	task, _ := queue.Pop()

	r.Equal(1, task.ID)
	r.Equal(0, queue.Len())
}

func TestQueueAutonomousTaskMustBeReady(t *testing.T) {
	r := require.New(t)

	queue := New()
	queue.Add(tasks.Task{
		ID: 1,
		Action: tasks.Command{
			Text: "echo test",
		},
	})
	task, _ := queue.Pop()

	r.Equal(status.Ready, task.Status)
}

func TestQueueDependentTaskMustBeWaiting(t *testing.T) {
	r := require.New(t)

	queue := New()
	queue.Add(tasks.Task{
		ID:      2,
		Depends: []int{1},
		Action: tasks.Command{
			Text: "true",
		},
	})
	task, _ := queue.Pop()

	r.Equal(status.Waiting, task.Status)
}

func TestQueueDependentTaskMustBeReady(t *testing.T) {
	r := require.New(t)

	queue := New()
	queue.Add(tasks.Task{
		ID: 1,
		Action: tasks.Command{
			Text: "true",
		},
	})
	queue.Add(tasks.Task{
		ID:      2,
		Depends: []int{1},
		Action: tasks.Command{
			Text: "true",
		},
	})

	_, _ = queue.Pop()
	queue.SetStatus(1, 0, status.Succeeded)
	task, _ := queue.Pop()

	r.Equal(status.Ready, task.Status)
}

func TestQueueDependentTaskMustBeInterrupted(t *testing.T) {
	r := require.New(t)

	queue := New()
	queue.Add(tasks.Task{
		ID: 1,
		Action: tasks.Command{
			Text: "true",
		},
	})
	queue.Add(tasks.Task{
		ID:      2,
		Depends: []int{1},
		Action: tasks.Command{
			Text: "true",
		},
	})

	_, _ = queue.Pop()
	queue.SetStatus(1, 0, status.Interrupted)
	task, _ := queue.Pop()

	r.Equal(status.Interrupted, task.Status)
}
