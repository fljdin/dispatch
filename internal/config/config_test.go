package config_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/config"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigWithDefaultMaxWorkers(t *testing.T) {
	config, _ := NewBuilder().Build()

	assert.Equal(t, 2, config.MaxWorkers)
}

func TestConfigWithMaxWorkers(t *testing.T) {
	config, _ := NewBuilder().
		WithMaxWorkers(1).
		Build()

	assert.Equal(t, 1, config.MaxWorkers)
}

func TestConfigWithMaxWorkersOverrided(t *testing.T) {
	yamlConfig := "workers: 2"

	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		WithMaxWorkers(1).
		Build()

	assert.Equal(t, 1, config.MaxWorkers)
}

func TestConfigFromYAML(t *testing.T) {
	yamlConfig := `
workers: 1
tasks:
  - id: 1
    command: true`

	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.GetTasks()

	assert.Equal(t, 1, config.MaxWorkers)
	assert.Equal(t, 1, tasks[0].ID)
}

func TestConfigFromYAMLWithDefaultConnection(t *testing.T) {
	yamlConfig := `
connections:
  - name: default
    host: remote
tasks:
  - id: 1
    type: psql
    command: SELECT 1
`
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.GetTasks()

	assert.Equal(t, 1, len(config.Connections))
	assert.Equal(t, "postgresql://?host=remote", tasks[0].Command.URI)
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
	config, _ := NewBuilder().
		WithYAML(fmt.Sprintf(yamlConfig, cnx)).
		Build()
	tasks, _ := config.GetTasks()

	assert.Equal(t, cnx, tasks[0].Command.URI)
}

func TestConfigFromYAMLWithUnknownConnection(t *testing.T) {
	yamlConfig := `
tasks:
  - id: 1
    command: SELECT 1
    connection: db
`
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	_, err := config.GetTasks()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "connection not found")
}

func TestConfigFromYAMLWithCombinedURIConnection(t *testing.T) {
	yamlConfig := `
connections:
  - name: db
    host: localhost
    port: 5433
    dbname: db
tasks:
  - id: 1
    type: psql
    command: SELECT 1
    connection: db
`

	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.GetTasks()

	expected := "postgresql://?dbname=db&host=localhost&port=5433"
	assert.Equal(t, expected, tasks[0].Command.URI)
}

func TestConfigFromNonExistingFile(t *testing.T) {
	yamlFilename := "test.yaml"

	_, err := NewBuilder().
		FromYAML(yamlFilename).
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestConfigFromInvalidYAML(t *testing.T) {
	yamlFilename := "config_*.yaml"
	yamlContent := "<xml></xml>"
	tempFile, _ := os.CreateTemp("", yamlFilename)

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	tempFile.Write([]byte(yamlContent))

	_, err := NewBuilder().
		FromYAML(tempFile.Name()).
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal")
}

func TestConfigLoadTasksFromFile(t *testing.T) {
	sqlFilename := "queries_*.sql"
	sqlContent := "SELECT 1; SELECT 2;"
	tempFile, _ := os.CreateTemp("", sqlFilename)

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	tempFile.Write([]byte(sqlContent))

	config, _ := NewBuilder().
		WithTask(YamlTask{
			ID:   1,
			Type: "psql",
			File: tempFile.Name(),
			URI:  "postgresql://localhost",
		}).
		Build()
	tasks, _ := config.GetTasks()

	// File task must be replaced by Command tasks loaded from SQL file
	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "SELECT 1;", tasks[0].Command.Text)
	assert.Equal(t, "postgresql://localhost", tasks[0].Command.URI)

	// Each loaded task must have an unique query ID
	assert.Equal(t, 0, tasks[0].QueryID)
	assert.Equal(t, 1, tasks[1].QueryID)
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
	_, err := NewBuilder().
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
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	_, err := config.GetTasks()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "depends on unknown task")
}

func TestConfigWithDefaultConnection(t *testing.T) {
	yamlConfig := `
connections:
  - name: anotherdb
    uri: postgresql://?host=localhost
tasks:
  - id: 1
    type: psql
    command: SELECT 1
  - id: 2
    type: psql
    command: SELECT 2
    connection: anotherdb
`
	cnx := Connection{Host: "remote", User: "postgres"}
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		WithDefaultConnection(cnx).
		Build()
	tasks, _ := config.GetTasks()

	assert.Equal(t, "postgresql://?host=remote&user=postgres", tasks[0].Command.URI)
	assert.Equal(t, "postgresql://?host=localhost", tasks[1].Command.URI)
}
