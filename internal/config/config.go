package config

import (
	"fmt"
	"runtime"

	"github.com/fljdin/dispatch/internal/parser"
	"github.com/fljdin/dispatch/internal/tasks"
)

var ConfigWorkersDefault int = 2

type YamlTask struct {
	ID         int    `yaml:"id"`
	Type       string `yaml:"type,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Command    string `yaml:"command"`
	File       string `yaml:"file"`
	URI        string `yaml:"uri,omitempty"`
	Connection string `yaml:"connection,omitempty"`
	Depends    []int  `yaml:"depends_on,omitempty"`
	ExecOutput string `yaml:"exec_output,omitempty"`
}

func (t YamlTask) Normalize() tasks.Task {
	return tasks.Task{
		ID:   t.ID,
		Name: t.Name,
		Command: tasks.Command{
			Text:       t.Command,
			File:       t.File,
			Type:       t.Type,
			URI:        t.URI,
			Connection: t.Connection,
			ExecOutput: t.ExecOutput,
		},
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

		// parse queries from file and append new related tasks
		if task.Command.File != "" {
			parser, err := parser.NewBuilder(task.Command.Type).
				FromFile(task.Command.File).
				Build()

			if err != nil {
				return nil, err
			}

			for queryId, query := range parser.Parse() {
				finalTasks = append(finalTasks, tasks.Task{
					ID:      task.ID,
					QueryID: queryId,
					Name:    fmt.Sprintf("Query loaded from %s", task.Command.File),
					Command: tasks.Command{
						Text: query,
						Type: task.Command.Type,
						URI:  task.Command.URI,
					},
				})
			}
		}

		// append task to already knwown identifiers
		identifiers = append(identifiers, task.ID)
	}

	return finalTasks, nil
}
