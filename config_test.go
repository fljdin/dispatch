package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfigFromBuilder(t *testing.T) {
	config, _ := NewConfigBuilder().
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

	config, _ := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	assert.Equal(t, config.Tasks[0].ID, 1)
}

func TestCreateConfigFromNonExistingFile(t *testing.T) {
	yamlFilename := "test.yaml"

	_, err := NewConfigBuilder().
		FromYAML(yamlFilename).
		Build()

	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestCreateConfigFromInvalidYAML(t *testing.T) {
	yamlFilename := "config_*.yaml"
	yamlContent := "<xml></xml>"
	tempFile, _ := os.CreateTemp("", yamlFilename)

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	tempFile.Write([]byte(yamlContent))

	_, err := NewConfigBuilder().
		FromYAML(tempFile.Name()).
		Build()

	assert.Contains(t, err.Error(), "cannot unmarshal")
}
