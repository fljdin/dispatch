package config_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

func TestConfigWithDefaultMaxWorkers(t *testing.T) {
	r := require.New(t)

	config, _ := NewBuilder().Build()

	r.Equal(2, config.MaxWorkers)
}

func TestConfigWithMaxWorkers(t *testing.T) {
	r := require.New(t)

	config, _ := NewBuilder().
		WithMaxWorkers(1).
		Build()

	r.Equal(1, config.MaxWorkers)
}

func TestConfigWithMaxWorkersOverrided(t *testing.T) {
	r := require.New(t)

	yamlConfig := "workers: 2"
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		WithMaxWorkers(1).
		Build()

	r.Equal(1, config.MaxWorkers)
}

func TestConfigFromYAML(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	workers: 1
	tasks:
	  - id: 1
	    command: true
	`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.Tasks()

	r.Equal(1, config.MaxWorkers)
	r.Equal(1, tasks[0].ID)
}

func TestConfigFromYAMLWithDefaultConnection(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	connections:
	  - name: default
	    host: remote
	tasks:
	  - id: 1
	    type: psql
	    command: SELECT 1`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.Tasks()

	r.Equal(1, len(config.Connections))
	r.Equal("postgresql://?host=remote", tasks[0].Action.(actions.Command).URI)
}

func TestConfigFromYAMLWithConnections(t *testing.T) {
	r := require.New(t)

	cnx := "postgresql://postgres:secret@localhost:5432/postgres"
	yamlConfig := dedent.Dedent(`
	connections:
	  - name: db
	    uri: %s
	tasks:
	  - id: 1
	    name: use predefined db name
	    type: psql
	    command: SELECT 1
	    connection: db`)
	config, _ := NewBuilder().
		WithYAML(fmt.Sprintf(yamlConfig, cnx)).
		Build()
	tasks, _ := config.Tasks()

	r.Equal(cnx, tasks[0].Action.(actions.Command).URI)
}

func TestConfigFromYAMLWithUnknownConnection(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    command: SELECT 1
	    connection: db`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	_, err := config.Tasks()

	r.NotNil(err)
	r.Contains(err.Error(), "connection not found")
}

func TestConfigFromYAMLWithCombinedURIConnection(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	connections:
	  - name: db
	    host: localhost
	    port: 5433
	    dbname: db
	tasks:
	  - id: 1
	    type: psql
	    command: SELECT 1
	    connection: db`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.Tasks()
	expected := "postgresql://?dbname=db&host=localhost&port=5433"

	r.Equal(expected, tasks[0].Action.(actions.Command).URI)
}

func TestConfigFromNonExistingFile(t *testing.T) {
	r := require.New(t)

	yamlFilename := "test.yaml"
	_, err := NewBuilder().
		FromYAML(yamlFilename).
		Build()

	r.NotNil(err)
	r.Contains(err.Error(), "no such file or directory")
}

func TestConfigFromInvalidYAML(t *testing.T) {
	r := require.New(t)

	yamlFilename := "config_*.yaml"
	yamlContent := "<xml></xml>"
	tempFile, _ := os.CreateTemp("", yamlFilename)
	tempFile.Write([]byte(yamlContent))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	_, err := NewBuilder().
		FromYAML(tempFile.Name()).
		Build()

	r.NotNil(err)
	r.Contains(err.Error(), "cannot unmarshal")
}
func TestConfigWithDependencies(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    command: true
	  - id: 2
	    command: true
	    depends_on: [1]`)
	_, err := NewBuilder().
		WithYAML(yamlConfig).
		Build()

	r.Equal(nil, err)
}

func TestConfigWithUnknownDependency(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    command: true
	  - id: 2
	    command: true
	    depends_on: [1, 3]`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	_, err := config.Tasks()

	r.NotNil(err)
	r.Contains(err.Error(), "depends on unknown task")
}

func TestConfigWithDefaultConnection(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
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
	    connection: anotherdb`)
	cnx := Connection{Host: "remote", User: "postgres"}
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		WithDefaultConnection(cnx).
		Build()
	tasks, _ := config.Tasks()

	r.Equal("postgresql://?host=remote&user=postgres", tasks[0].Action.(actions.Command).URI)
	r.Equal("postgresql://?host=localhost", tasks[1].Action.(actions.Command).URI)
}

func TestConfigWithOutputLoader(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    loaded:
	      from: sh
	      command: echo true`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.Tasks()

	r.Equal("sh", tasks[0].Action.(actions.OutputLoader).From)
	r.Equal("echo true", tasks[0].Action.(actions.OutputLoader).Text)
}

func TestConfigWithFileLoader(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    type: psql
	    file: junk.sql`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, err := config.Tasks()

	r.Nil(err)
	r.Equal("junk.sql", tasks[0].Action.(actions.FileLoader).File)
}

func TestConfigWithInvalidLoader(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    loaded:
	      from: invalid`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	_, err := config.Tasks()

	r.NotNil(err)
	r.Equal("action is required", err.Error())
}
