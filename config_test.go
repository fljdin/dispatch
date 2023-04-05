package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigWithDefaultMaxWorkers(t *testing.T) {
	config, _ := NewConfigBuilder().Build()

	assert.Equal(t, 2, config.MaxWorkers)
}

func TestConfigWithMaxWorkers(t *testing.T) {
	config, _ := NewConfigBuilder().
		WithMaxWorkers(1).
		Build()

	assert.Equal(t, 1, config.MaxWorkers)
}

func TestConfigWithTask(t *testing.T) {
	config, _ := NewConfigBuilder().
		WithTask(Task{
			ID:      1,
			Command: "echo test",
		}).
		Build()

	assert.Equal(t, 1, config.Tasks[0].ID)
}

func TestConfigFromYAML(t *testing.T) {
	yamlConfig := `
workers: 1
tasks:
  - id: 1
    command: echo test`

	config, _ := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	assert.Equal(t, 1, config.MaxWorkers)
	assert.Equal(t, 1, config.Tasks[0].ID)
}

func TestConfigFromYAMLWithURIOnTask(t *testing.T) {
	cnx := "postgresql://postgres:secret@localhost:5432/postgres"
	yamlConfig := `
tasks:
  - id: 1
    name: use predefined db name
    type: psql
    command: SELECT 1
    uri: %s
`
	config, err := NewConfigBuilder().
		WithYAML(fmt.Sprintf(yamlConfig, cnx)).
		Build()

	assert.Equal(t, nil, err)
	assert.Equal(t, cnx, config.Tasks[0].URI)
}

func TestConfigFromYAMLWithConnections(t *testing.T) {
	cnx := "postgresql://postgres:secret@localhost:5432/postgres"
	yamlConfig := `
connections:
  - name: db
    uri: %s
tasks:
  - id: 1
    name: use predefined db name
    type: psql
    command: SELECT 1
    connection: db
`
	config, _ := NewConfigBuilder().
		WithYAML(fmt.Sprintf(yamlConfig, cnx)).
		Build()

	assert.Equal(t, 1, len(config.Connections))
	assert.Equal(t, cnx, config.Tasks[0].URI)
}

func TestConfigFromYAMLWithUnknownConnection(t *testing.T) {
	yamlConfig := `
tasks:
  - id: 1
    command: SELECT 1
    connection: db
`
	_, err := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	assert.Contains(t, err.Error(), "connection not found")
}

func TestConfigWithMaxWorkersOverrided(t *testing.T) {
	yamlConfig := "workers: 1"

	config, _ := NewConfigBuilder().
		WithYAML(yamlConfig).
		WithMaxWorkers(2).
		Build()

	assert.Equal(t, 2, config.MaxWorkers)
}

func TestConfigFromNonExistingFile(t *testing.T) {
	yamlFilename := "test.yaml"

	_, err := NewConfigBuilder().
		FromYAML(yamlFilename).
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "no such file or directory")
	}
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

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "cannot unmarshal")
	}
}

func TestConfigWithInvalidType(t *testing.T) {
	yamlConfig := `
tasks:
  - id: 1
    type: unknown
    command: unknown
`
	_, err := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "invalid task type")
	}
}

func TestConfigWithInvalidFileType(t *testing.T) {
	sqlFilename := "whatever.sql"
	tempFile, _ := os.CreateTemp("", sqlFilename)

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	yamlConfig := `
tasks:
  - id: 1
    type: sh
    file: %s
`
	_, err := NewConfigBuilder().
		WithYAML(fmt.Sprintf(yamlConfig, tempFile.Name())).
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "invalid type for parsing file")
	}
}
