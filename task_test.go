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

func TestShellTask(t *testing.T) {
	task := &Task{
		ID:      1,
		Command: "echo test",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Succeeded)
	assert.Contains(t, result.Message, "test")
}

func TestShellTaskWithError(t *testing.T) {
	task := &Task{
		ID:      1,
		Command: "false",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Failed)
}

func TestPsqlTask(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Execute SQL statement",
		Type:       "psql",
		Command:    "CREATE TEMPORARY TABLE foo (bar smallint);",
		Connection: "postgresql://postgres:secret@localhost:5432/postgres",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Succeeded)
	assert.Contains(t, result.Message, "CREATE TABLE")
}
