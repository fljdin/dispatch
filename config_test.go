package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfigFromBuilder(t *testing.T) {
	config := NewConfigBuilder().
		WithTask(Task{
			ID:      1,
			Command: "echo test",
		}).
		Build()

	assert.Equal(t, config.Tasks[0].ID, 1)
}

func TestCreateConfigFromYAML(t *testing.T) {
	yamlConfig := `
tasks:
  - id: 1
    command: echo test`

	config := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	assert.Equal(t, config.Tasks[0].ID, 1)
}
