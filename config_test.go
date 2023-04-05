package main

import (
	"fmt"
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

	assert.Equal(t, err, nil)
	assert.Equal(t, config.Tasks[0].URI, cnx)
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

	assert.Equal(t, len(config.Connections), 1)
	assert.Equal(t, config.Tasks[0].URI, cnx)
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
		WithMaxWorkers(4).
		Build()

	assert.Equal(t, config.MaxWorkers, 4)
}

func TestConfigFromNonExistingFile(t *testing.T) {
	yamlFilename := "test.yaml"

	_, err := NewConfigBuilder().
		FromYAML(yamlFilename).
		Build()

	if assert.NotEqual(t, err, nil) {
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

	if assert.NotEqual(t, err, nil) {
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

	if assert.NotEqual(t, err, nil) {
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

	if assert.NotEqual(t, err, nil) {
		assert.Contains(t, err.Error(), "invalid type for parsing file")
	}
}
