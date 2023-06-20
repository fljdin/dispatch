package dispatcher_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	. "github.com/fljdin/dispatch/internal/task"
	"github.com/stretchr/testify/assert"
)

func TestDispatcherAddTask(t *testing.T) {
	dispatcher, _ := NewDispatcherBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		Command: Command{Text: "true"},
	})

	assert.Equal(t, Waiting, dispatcher.GetStatus(1))
}

func TestDispatcherDependentTaskNeverExecuted(t *testing.T) {
	dispatcher, _ := NewDispatcherBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		Command: Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		ID:      2,
		Depends: []int{1},
		Command: Command{Text: "true"},
	})
	dispatcher.Wait()

	assert.Equal(t, Failed, dispatcher.GetStatus(1))
	assert.Equal(t, Interrupted, dispatcher.GetStatus(2))
}

func TestDispatcherDependentTaskGetSucceeded(t *testing.T) {
	dispatcher, _ := NewDispatcherBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		Command: Command{Text: "true"},
	})
	dispatcher.AddTask(Task{
		ID:      2,
		Depends: []int{1},
		Command: Command{Text: "true"},
	})
	dispatcher.Wait()

	assert.Equal(t, Succeeded, dispatcher.GetStatus(1))
	assert.Equal(t, Succeeded, dispatcher.GetStatus(2))
}

func TestDispatcherStatusOfFileTaskMustSummarizeLoadedTaskStatus(t *testing.T) {
	dispatcher, _ := NewDispatcherBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		QueryID: 0,
		Command: Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		ID:      1,
		QueryID: 1,
		Command: Command{Text: "true"},
	})
	dispatcher.Wait()

	assert.Equal(t, Failed, dispatcher.GetStatus(1))
}
