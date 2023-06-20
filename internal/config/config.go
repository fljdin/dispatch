package config

import (
	"runtime"

	"github.com/fljdin/dispatch/internal/tasks"
)

var ConfigWorkersDefault int = 2

type YamlGenerator struct {
	From    string `yaml:"from"`
	Command string `yaml:"command"`
	File    string `yaml:"file"`
}

type YamlTask struct {
	ID         int           `yaml:"id"`
	Type       string        `yaml:"type,omitempty"`
	Name       string        `yaml:"name,omitempty"`
	Command    string        `yaml:"command"`
	File       string        `yaml:"file"`
	URI        string        `yaml:"uri,omitempty"`
	Connection string        `yaml:"connection,omitempty"`
	Depends    []int         `yaml:"depends_on,omitempty"`
	Generated  YamlGenerator `yaml:"generated,omitempty"`
}

func (t YamlTask) Normalize() tasks.Task {
	command := tasks.Command{
		Text:       t.Command,
		File:       t.File,
		Type:       t.Type,
		URI:        t.URI,
		Connection: t.Connection,
	}

	if t.Generated.From != "" {
		command = tasks.Command{
			Text:       t.Generated.Command,
			File:       t.Generated.File,
			Type:       t.Type,
			URI:        t.URI,
			Connection: t.Connection,
			From:       t.Generated.From,
		}
	}

	return tasks.Task{
		ID:      t.ID,
		Name:    t.Name,
		Command: command,
		Depends: t.Depends,
	}
}

type Config struct {
	DeclaredTasks     []YamlTask        `yaml:"tasks"`
	MaxWorkers        int               `yaml:"workers"`
	Logfile           string            `yaml:"logfile"`
	Connections       tasks.Connections `yaml:"connections"`
	DefaultConnection tasks.Connection
}

func (c *Config) ConfigureWorkers() {
	if c.MaxWorkers < 1 {
		c.MaxWorkers = ConfigWorkersDefault
	}

	if c.MaxWorkers > runtime.NumCPU() {
		c.MaxWorkers = runtime.NumCPU()
	}
}

func (c *Config) ConfigureConnections() {
	if _, err := c.Connections.GetURIByName("default"); err != nil {
		c.Connections = append(c.Connections, tasks.Connection{
			Name: "default",
			URI:  c.DefaultConnection.CombinedURI(),
		})
	}
}

func (c Config) Tasks() ([]tasks.Task, error) {
	var finalTasks []tasks.Task
	var identifiers []int

	for _, declared := range c.DeclaredTasks {
		task := declared.Normalize()

		if err := task.VerifyRequired(); err != nil {
			return nil, err
		}

		if err := task.Command.VerifyType(); err != nil {
			return nil, err
		}

		if err := task.VerifyDependencies(identifiers); err != nil {
			return nil, err
		}

		// auto-complete URI from named connections
		if task.Command.Connection != "" {
			uri, err := c.Connections.GetURIByName(task.Command.Connection)

			if err != nil {
				return nil, err
			}

			task.Command.URI = uri
		}

		// use default connection if no URI is provided
		if task.Command.URI == "" {
			task.Command.URI, _ = c.Connections.GetURIByName("default")
		}

		// append task to final tasks
		if task.Command.Text != "" {
			finalTasks = append(finalTasks, task)
		}

		// append task to already knwown identifiers
		identifiers = append(identifiers, task.ID)
	}

	return finalTasks, nil
}
