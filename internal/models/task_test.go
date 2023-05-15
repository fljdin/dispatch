package models_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	task := Task{
		ID:      1,
		Command: "echo test",
	}
	err := task.VerifyRequired()

	assert.Equal(t, nil, err)
}

func TestTaskVerifyIDRequired(t *testing.T) {
	task := Task{
		Command: "true",
	}
	err := task.VerifyRequired()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "id is required")
}

func TestTaskVerifyCommandRequired(t *testing.T) {
	task := Task{
		ID: 1,
	}
	err := task.VerifyRequired()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "command is required")
}

func TestTaskVerifyType(t *testing.T) {
	task := Task{
		ID:      1,
		Type:    "unknown",
		Command: "unknown",
	}
	err := task.VerifyType()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid task type")
}

func TestShellTaskWithOutput(t *testing.T) {
	task := Task{
		ID:      1,
		Command: "echo test",
	}
	result := task.Run()

	assert.Equal(t, Succeeded, result.Status)
	assert.Contains(t, result.Output, "test")
}

func TestShellTaskWithError(t *testing.T) {
	task := Task{
		ID:      1,
		Command: "false",
	}
	result := task.Run()

	assert.Equal(t, Failed, result.Status)
}

func TestWriteOutputToFile(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "task_*.out")

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	task := Task{
		ID:      1,
		Command: "echo test",
		Output:  tempFile.Name(),
	}
	task.Run()

	data, err := os.ReadFile(tempFile.Name())
	if assert.Equal(t, nil, err) {
		assert.Equal(t, "test\n", string(data))
	}
}

func TestTaskOutputWithInvalidTemplating(t *testing.T) {
	task := Task{
		ID:      1,
		Command: "true",
		Output:  "task_{{with .Invalid}}.out",
	}
	err := task.VerifyOutput()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "unexpected")
}

func TestTaskOutputWithTemplating(t *testing.T) {
	task := Task{
		ID:      1,
		QueryID: 0,
		Command: "true",
		Output:  "task_{{.ID}}_{{.QueryID}}.out",
	}
	err := task.VerifyOutput()

	require.Nil(t, err)
	assert.Equal(t, "task_1_0.out", task.GetOutput())
}
