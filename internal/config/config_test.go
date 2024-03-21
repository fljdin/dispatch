package config_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/config"
	. "github.com/fljdin/dispatch/internal/tasks"
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

func TestConfigFromYAMLWithDefaultEnvironment(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	environments:
	  - name: default
	    variables:
	      key: bar
	tasks:
	  - id: 1
	    command: true`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.Tasks()

	r.Equal("bar", tasks[0].Action.(Command).Variables["key"])
}

func TestConfigFromYAMLWithEnvironment(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	environments:
	  - name: custom
	    variables:
	      key: foo
	  - name: default
	    variables:
	      key: bar
	tasks:
	  - id: 1
	    command: true
	    env: custom`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.Tasks()

	r.Equal("foo", tasks[0].Action.(Command).Variables["key"])
}

func TestConfigFromYAMLWithUnknownEnvironment(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    command: true
	    env: custom`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	_, err := config.Tasks()

	r.NotNil(err)
	r.Contains(err.Error(), "environment not found")
}

func TestConfigFromYAMLWithTaskVariables(t *testing.T) {
	r := require.New(t)

	yamlConfig := dedent.Dedent(`
	tasks:
	  - id: 1
	    command: true
	    variables:
	      key: bar`)
	config, _ := NewBuilder().
		WithYAML(yamlConfig).
		Build()
	tasks, _ := config.Tasks()

	r.Equal("bar", tasks[0].Action.(Command).Variables["key"])
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

	r.Equal("sh", tasks[0].Action.(OutputLoader).From)
	r.Equal("echo true", tasks[0].Action.(OutputLoader).Text)
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
	r.Equal("junk.sql", tasks[0].Action.(FileLoader).File)
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
