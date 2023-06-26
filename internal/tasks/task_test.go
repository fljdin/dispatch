package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	task := Task{
		ID: 1,
		Action: actions.Command{
			Text: "echo test",
		},
	}
	err := task.Validate()

	assert.Equal(t, nil, err)
}

func TestTaskVerifyIDRequired(t *testing.T) {
	task := Task{
		Action: actions.Command{Text: "true"},
	}
	err := task.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "id is required")
}

func TestTaskVerifyCommandRequired(t *testing.T) {
	task := Task{
		ID: 1,
	}
	err := task.Validate()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "action is required")
}
