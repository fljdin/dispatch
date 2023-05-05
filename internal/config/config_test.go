package config_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/config"
	. "github.com/fljdin/dispatch/internal/models"
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

func TestConfigWithMaxWorkersOverrided(t *testing.T) {
	yamlConfig := "workers: 2"

	config, _ := NewConfigBuilder().
		WithYAML(yamlConfig).
		WithMaxWorkers(1).
		Build()

	assert.Equal(t, 1, config.MaxWorkers)
}

func TestConfigWithTask(t *testing.T) {
	config, _ := NewConfigBuilder().
		WithTask(Task{
			ID:      1,
			Command: "true",
		}).
		Build()

	assert.Equal(t, 1, config.Tasks[0].ID)
}

func TestConfigFromYAML(t *testing.T) {
	yamlConfig := `
workers: 1
tasks:
  - id: 1
    command: true`

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

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "connection not found")
	}
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

func TestConfigLoadTasksFromFile(t *testing.T) {
	sqlFilename := "queries_*.sql"
	sqlContent := "SELECT 1; SELECT 2;"
	tempFile, _ := os.CreateTemp("", sqlFilename)

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	tempFile.Write([]byte(sqlContent))

	config, _ := NewConfigBuilder().
		WithTask(Task{
			ID:   1,
			Type: "psql",
			File: tempFile.Name(),
			URI:  "postgresql://localhost",
		}).
		Build()

	// File task must be replaced by Command tasks loaded from SQL file
	assert.Equal(t, 1, config.Tasks[0].ID)
	assert.Equal(t, "SELECT 1;", config.Tasks[0].Command)
	assert.Equal(t, "postgresql://localhost", config.Tasks[0].URI)

	// Each loaded task must have an unique query ID
	assert.Equal(t, 0, config.Tasks[0].QueryID)
	assert.Equal(t, 1, config.Tasks[1].QueryID)
}

func TestConfigWithDependencies(t *testing.T) {
	yamlConfig := `
tasks:
  - id: 1
    command: true
  - id: 2
    command: true
    depends_on: [1]
`
	_, err := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	assert.Equal(t, nil, err)
}

func TestConfigWithUnknownDependency(t *testing.T) {
	yamlConfig := `
tasks:
  - id: 1
    command: true
  - id: 2
    command: true
    depends_on: [1, 3]
`
	_, err := NewConfigBuilder().
		WithYAML(yamlConfig).
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "depends on unknown task")
	}
}