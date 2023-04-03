package main

import (
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

func TestRunShellTaskShouldSucceed(t *testing.T) {
	results := make(chan TaskResult, 1)
	task := &Task{
		ID:      1,
		Command: "echo test",
	}

	task.Run(results)
	result := <-results

	assert.Equal(t, result.Status, Succeeded)
	assert.Equal(t, result.Message, "test\n")
}
