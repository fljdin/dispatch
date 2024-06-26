package routines_test

import (
	"testing"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/routines"
	"github.com/fljdin/dispatch/internal/status"
	"github.com/stretchr/testify/require"
)

func TestQueueNext(t *testing.T) {
	r := require.New(t)
	var (
		ok    bool
		ready config.Task
	)

	queue := routines.NewQueue()
	queue.Add(config.Task{
		Identifier: config.NewId(1, 0),
		Status:     status.Waiting,
	})

	ready, ok = queue.Next()
	r.True(ok) // task 1 is immediately ready
	r.Equal(1, ready.Identifier.ID)

	queue.Update(config.NewId(1, 0), status.Succeeded)

	_, ok = queue.Next()
	r.False(ok) // queue is empty
}

func TestQueueEvaluate(t *testing.T) {
	r := require.New(t)

	queue := routines.NewQueue()
	queue.Add(config.Task{Identifier: config.NewId(1, 0), Status: status.Succeeded})
	queue.Add(config.Task{Identifier: config.NewId(1, 1), Status: status.Failed})

	r.Equal(status.Failed, queue.Evaluate(1))

	queue.Add(config.Task{Identifier: config.NewId(2, 0), Status: status.Succeeded})
	queue.Add(config.Task{Identifier: config.NewId(2, 1), Status: status.Waiting})

	r.Equal(status.Waiting, queue.Evaluate(2))

	queue.Add(config.Task{Identifier: config.NewId(3, 0), Status: status.Succeeded})
	queue.Add(config.Task{Identifier: config.NewId(3, 1), Status: status.Succeeded})

	r.Equal(status.Succeeded, queue.Evaluate(3))

	queue.Add(config.Task{Identifier: config.NewId(4, 0), Status: status.Succeeded})
	queue.Add(config.Task{Identifier: config.NewId(4, 1), Status: status.Interrupted})

	r.Equal(status.Interrupted, queue.Evaluate(4))
}

func TestQueueTaskWithDependencies(t *testing.T) {
	r := require.New(t)
	var (
		ok    bool
		ready config.Task
	)

	queue := routines.NewQueue()
	queue.Add(config.Task{
		Identifier: config.NewId(1, 0),
		Status:     status.Waiting,
		Depends:    []int{2},
	})

	_, ok = queue.Next()
	r.False(ok) // task 1 is waiting for task 2

	queue.Add(config.Task{
		Identifier: config.NewId(2, 0),
		Status:     status.Waiting,
	})

	ready, ok = queue.Next()
	r.True(ok) // task 2 is ready
	r.Equal(2, ready.Identifier.ID)

	queue.Update(config.NewId(2, 0), status.Succeeded)

	ready, ok = queue.Next()
	r.True(ok) // task 1 is ready
	r.Equal(1, ready.Identifier.ID)
}

func TestQueueTaskIsInterrupted(t *testing.T) {
	r := require.New(t)
	var (
		ok    bool
		ready config.Task
	)

	queue := routines.NewQueue()
	queue.Add(config.Task{
		Identifier: config.NewId(1, 0),
		Status:     status.Waiting,
		Depends:    []int{2},
	})

	ready, ok = queue.Next()
	r.False(ok) // task 1 is waiting for task 2

	queue.Add(config.Task{
		Identifier: config.NewId(2, 0),
		Status:     status.Waiting,
	})
	queue.Update(config.NewId(2, 0), status.Succeeded)

	queue.Add(config.Task{
		Identifier: config.NewId(2, 1),
		Status:     status.Waiting,
	})
	queue.Update(config.NewId(2, 1), status.Failed)

	ready, ok = queue.Next()
	r.True(ok) // task 1 is ready
	r.Equal(1, ready.Identifier.ID)
	r.Equal(status.Interrupted, ready.Status)
}
