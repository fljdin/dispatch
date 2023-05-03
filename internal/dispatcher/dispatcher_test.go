package dispatcher_test

import (
	"context"
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	. "github.com/fljdin/dispatch/internal/models"
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

func TestDispatcherTraceToFile(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "trace_*.out")

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	dispatcher := NewDispatcher(context.Background(), 1, 1)
	dispatcher.TraceTo(tempFile.Name())

	dispatcher.Add(Task{
		ID:      1,
		Command: "echo test",
	})
	dispatcher.Wait()

	data, err := os.ReadFile(tempFile.Name())
	if assert.Equal(t, nil, err) {
		assert.Contains(t, string(data), "test\n")
	}
}
