package config_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/actions"
	. "github.com/fljdin/dispatch/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	r := require.New(t)

	task := Task{
		Identifier: NewId(1, 0),
		Action: Command{
			Type: Shell,
			Text: "echo test",
		},
	}
	err := task.Validate()

	r.Equal(nil, err)
}

func TestTaskVerifyIDRequired(t *testing.T) {
	r := require.New(t)

	task := Task{
		Action: Command{Text: "true"},
	}
	err := task.Validate()

	r.NotNil(err)
	r.Contains(err.Error(), "id is required")
}

func TestTaskVerifyCommandRequired(t *testing.T) {
	r := require.New(t)

	task := Task{
		Identifier: NewId(1, 0),
	}
	err := task.Validate()

	r.NotNil(err)
	r.Contains(err.Error(), "action is required")
}
