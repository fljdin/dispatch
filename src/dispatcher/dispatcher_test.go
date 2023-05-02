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

	assert.Equal(t, Succeeded, dispatcher.GetStatus(1))
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

	assert.Equal(t, Failed, dispatcher.GetStatus(1))
	assert.Equal(t, Interrupted, dispatcher.GetStatus(2))
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

	assert.Equal(t, Succeeded, dispatcher.GetStatus(1))
	assert.Equal(t, Succeeded, dispatcher.GetStatus(2))
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

	assert.Equal(t, Failed, dispatcher.GetStatus(1))
}
