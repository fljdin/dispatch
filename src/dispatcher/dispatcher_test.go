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

	result := <-dispatcher.Results
	assert.Equal(t, Succeeded, result.Status)
}
