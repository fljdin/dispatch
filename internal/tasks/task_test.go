package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	r := require.New(t)

	task := Task{
		ID: 1,
		Action: actions.Command{
			Text: "echo test",
		},
	}
	err := task.Validate()

	r.Equal(nil, err)
}

func TestTaskVerifyIDRequired(t *testing.T) {
	r := require.New(t)

	task := Task{
		Action: actions.Command{Text: "true"},
	}
	err := task.Validate()

	r.NotNil(err)
	r.Contains(err.Error(), "id is required")
}

func TestTaskVerifyCommandRequired(t *testing.T) {
	r := require.New(t)

	task := Task{
		ID: 1,
	}
	err := task.Validate()

	r.NotNil(err)
	r.Contains(err.Error(), "action is required")
}
