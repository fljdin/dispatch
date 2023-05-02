package dispatcher_test

import (
	"context"
	"testing"

	. "github.com/fljdin/dispatch/src/dispatcher"
	. "github.com/fljdin/dispatch/src/models"
	"github.com/stretchr/testify/assert"
)

func TestDispatcherAddTask(t *testing.T) {
	dispatcher := NewDispatcher(context.Background(), 1, 1)
	dispatcher.Add(Task{
		ID:      1,
		Command: "true",
	})
	dispatcher.Wait()

	status, ok := dispatcher.GetStatus(1)
	if assert.Equal(t, true, ok) {
		assert.Equal(t, Succeeded, status)
	}
}

func TestDispatcherDependentTaskNeverExecuted(t *testing.T) {
	dispatcher := NewDispatcher(context.Background(), 1, 2)
	dispatcher.Add(Task{
		ID:      1,
		Command: "false",
	})
	dispatcher.Add(Task{
		ID:      2,
		Depends: []int{1},
		Command: "true",
	})
	dispatcher.Wait()

	status, ok := dispatcher.GetStatus(1)
	if assert.Equal(t, true, ok) {
		assert.Equal(t, Failed, status)
	}

	status, ok = dispatcher.GetStatus(2)
	if assert.Equal(t, true, ok) {
		assert.Equal(t, Interrupted, status)
	}
}

func TestDispatcherDependentTaskGetSucceeded(t *testing.T) {
	dispatcher := NewDispatcher(context.Background(), 1, 2)
	dispatcher.Add(Task{
		ID:      1,
		Command: "true",
	})
	dispatcher.Add(Task{
		ID:      2,
		Depends: []int{1},
		Command: "true",
	})
	dispatcher.Wait()

	status, ok := dispatcher.GetStatus(1)
	if assert.Equal(t, true, ok) {
		assert.Equal(t, Succeeded, status)
	}

	status, ok = dispatcher.GetStatus(2)
	if assert.Equal(t, true, ok) {
		assert.Equal(t, Succeeded, status)
	}
}

func TestDispatcherStatusOfFileTaskMustSummarizeLoadedTaskStatus(t *testing.T) {
	dispatcher := NewDispatcher(context.Background(), 1, 2)
	dispatcher.Add(Task{
		ID:      1,
		QueryID: 0,
		Command: "false",
	})
	dispatcher.Add(Task{
		ID:      1,
		QueryID: 1,
		Command: "true",
	})
	dispatcher.Wait()

	status, _ := dispatcher.GetStatus(1)
	assert.Equal(t, Failed, status)
}
