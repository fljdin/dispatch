package config_test

import (
	"os"
	"testing"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

func TestConfigFromYAML(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  command: true
	`)
	cfg, _ := config.NewFromRaw(yaml)
	tasks, _ := cfg.Tasks()

	r.Equal(1, tasks[0].Identifier.ID)
}

func TestConfigFromYAMLWithDefaultEnvironment(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		environments:
		- name: default
		  variables:
		    key: bar
		tasks:
		- id: 1
		  command: true
	`)
	cfg, _ := config.NewFromRaw(yaml)
	ts, _ := cfg.Tasks()

	r.Equal("bar", ts[0].Action.(tasks.Command).Variables["key"])
}

func TestConfigFromYAMLWithEnvironment(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
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
		  env: custom
	`)
	cfg, _ := config.NewFromRaw(yaml)
	ts, _ := cfg.Tasks()

	r.Equal("foo", ts[0].Action.(tasks.Command).Variables["key"])
}

func TestConfigFromYAMLWithUnknownEnvironment(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  command: true
		  env: custom
	`)
	cfg, _ := config.NewFromRaw(yaml)
	_, err := cfg.Tasks()

	r.Error(err)
	r.Contains(err.Error(), "environment not found")
}

func TestConfigFromYAMLWithTaskVariables(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  command: true
		  variables:
		    key: bar
	`)
	cfg, _ := config.NewFromRaw(yaml)
	ts, _ := cfg.Tasks()

	r.Equal("bar", ts[0].Action.(tasks.Command).Variables["key"])
}

func TestConfigFromNonExistingFile(t *testing.T) {
	r := require.New(t)

	path := "test.yaml"
	_, err := config.New(path)

	r.Error(err)
	r.Contains(err.Error(), "no such file or directory")
}

func TestConfigFromInvalidYAML(t *testing.T) {
	r := require.New(t)

	path := "config_*.yaml"
	yaml := "<xml></xml>"
	tempFile, _ := os.CreateTemp("", path)
	tempFile.Write([]byte(yaml))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	_, err := config.New(tempFile.Name())

	r.Error(err)
	r.Contains(err.Error(), "cannot unmarshal")
}
func TestConfigWithDependencies(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  command: true
		- id: 2
		  command: true
		  depends_on: [1]
	`)
	_, err := config.NewFromRaw(yaml)

	r.NoError(err)
}

func TestConfigWithUnknownDependency(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  command: true
		- id: 2
		  command: true
		  depends_on: [1, 3]
	`)
	cfg, _ := config.NewFromRaw(yaml)
	_, err := cfg.Tasks()

	r.Error(err)
	r.Contains(err.Error(), "depends on unknown task")
}

func TestConfigWithOutputLoader(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  loaded:
		    from: sh
		    command: echo true
	`)
	cfg, _ := config.NewFromRaw(yaml)
	ts, _ := cfg.Tasks()

	r.Equal("sh", ts[0].Action.(tasks.OutputLoader).From)
	r.Equal("echo true", ts[0].Action.(tasks.OutputLoader).Text)
}

func TestConfigWithFileLoader(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  type: psql
		  file: junk.sql
	`)
	cfg, _ := config.NewFromRaw(yaml)
	ts, err := cfg.Tasks()

	r.NoError(err)
	r.Equal("junk.sql", ts[0].Action.(tasks.FileLoader).File)
}

func TestConfigWithInvalidLoader(t *testing.T) {
	r := require.New(t)

	yaml := dedent.Dedent(`
		tasks:
		- id: 1
		  loaded:
		    from: invalid
	`)
	cfg, _ := config.NewFromRaw(yaml)
	_, err := cfg.Tasks()

	r.Error(err)
	r.Equal("action is required", err.Error())
}
