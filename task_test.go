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

func TestTaskShouldSucceed(t *testing.T) {
	task := &Task{
		ID:      1,
		Command: "echo test",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Succeeded)
	assert.Equal(t, result.Message, "test\n")
}

func TestTaskShouldFail(t *testing.T) {
	task := &Task{
		ID:      1,
		Command: "false",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Failed)
}
