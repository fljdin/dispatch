package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	task := &Task{
		ID:      1,
		Command: "echo test",
	}

	assert.Equal(t, task.ID, 1)
}

func TestTaskVerifyIDRequired(t *testing.T) {
	task := &Task{
		Command: "true",
	}
	err := task.VerifyRequired()

	if assert.NotEqual(t, err, nil) {
		assert.Contains(t, err.Error(), "id is required")
	}
}

func TestTaskVerifyCommandRequired(t *testing.T) {
	task := &Task{
		ID: 1,
	}
	err := task.VerifyRequired()

	if assert.NotEqual(t, err, nil) {
		assert.Contains(t, err.Error(), "command is required")
	}
}

func TestTaskVerifyType(t *testing.T) {
	task := &Task{
		ID:      1,
		Type:    "unknown",
		Command: "unknown",
	}
	err := task.VerifyType()

	if assert.NotEqual(t, err, nil) {
		assert.Contains(t, err.Error(), "invalid task type")
	}
}

func TestShellTask(t *testing.T) {
	task := &Task{
		ID:      1,
		Command: "echo test",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Succeeded)
	assert.Contains(t, result.Output, "test")
}

func TestShellTaskWithError(t *testing.T) {
	task := &Task{
		ID:      1,
		Command: "false",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Failed)
}
