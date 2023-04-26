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

	result := dispatcher.GetResult(1)
	assert.Equal(t, Succeeded, result.Status)
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

	result := dispatcher.GetResult(1)
	assert.Equal(t, Failed, result.Status)

	result = dispatcher.GetResult(2)
	assert.Equal(t, Waiting, result.Status)
}
