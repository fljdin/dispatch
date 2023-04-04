package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigWithDefaultMaxWorkers(t *testing.T) {
	config, _ := NewConfigBuilder().Build()

	assert.Equal(t, config.MaxWorkers, 2)
}

func TestConfigWithMaxWorkers(t *testing.T) {
	config, _ := NewConfigBuilder().
		WithMaxWorkers(4).
		Build()

	assert.Equal(t, config.MaxWorkers, 4)
}

func TestConfigWithTask(t *testing.T) {
	config, _ := NewConfigBuilder().
		WithTask(Task{
			ID:      1,
			Command: "echo test",
		}).
		Build()

	assert.Equal(t, config.Tasks[0].ID, 1)
}

func TestConfigFromYAML(t *testing.T) {
	yamlConfig := `
workers: 4
tasks:
  - id: 1
    command: echo test`

	config, _ := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	assert.Equal(t, config.MaxWorkers, 4)
	assert.Equal(t, config.Tasks[0].ID, 1)
}

func TestConfigWithMaxWorkersOverrided(t *testing.T) {
	yamlConfig := "workers: 1"

	config, _ := NewConfigBuilder().
		WithYAML(yamlConfig).
		WithMaxWorkers(4).
		Build()

	assert.Equal(t, config.MaxWorkers, 4)
}

func TestConfigFromNonExistingFile(t *testing.T) {
	yamlFilename := "test.yaml"

	_, err := NewConfigBuilder().
		FromYAML(yamlFilename).
		Build()

	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestConfigFromInvalidYAML(t *testing.T) {
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
