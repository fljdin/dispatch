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

func TestPsqlTask(t *testing.T) {
	task := &Task{
		ID:      1,
		Type:    "psql",
		Command: "CREATE TEMPORARY TABLE foo (bar smallint);",
		URI:     "postgresql://postgres:secret@localhost:5432/postgres",
	}
	result := task.Run(context.Background())

	assert.Equal(t, result.Status, Succeeded)
	assert.Contains(t, result.Output, "CREATE TABLE")
}

func TestPsqlTaskWihoutURI(t *testing.T) {
	task := &Task{
		ID:      1,
		Name:    "should connect to default socket",
		Type:    "psql",
		Command: "SELECT test",
	}
	result := task.Run(context.Background())

	assert.NotEqual(t, result.Status, Succeeded)
	assert.Contains(t, result.Output, "test")
}
