package config

import (
	"fmt"
	"runtime"

	"github.com/fljdin/dispatch/internal/parser"
	"github.com/fljdin/dispatch/internal/task"
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
}

type Config struct {
	DeclaredTasks     []YamlTask       `yaml:"tasks"`
	MaxWorkers        int              `yaml:"workers"`
	Logfile           string           `yaml:"logfile"`
	Connections       task.Connections `yaml:"connections"`
	DefaultConnection task.Connection
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
		c.Connections = append(c.Connections, task.Connection{
			Name: "default",
			URI:  c.DefaultConnection.CombinedURI(),
		})
	}
}

func (c Config) GetTasks() ([]task.Task, error) {
	var finalTasks []task.Task
	var identifiers []int

	for _, t := range c.DeclaredTasks {
		t := task.Task{
			ID:   t.ID,
			Name: t.Name,
			Command: task.Command{
				Text:       t.Command,
				File:       t.File,
				Type:       t.Type,
				URI:        t.URI,
				Connection: t.Connection,
			},
			Depends: t.Depends,
		}

		if err := t.VerifyRequired(); err != nil {
			return nil, err
		}

		if err := t.Command.VerifyType(); err != nil {
			return nil, err
		}

		if err := t.VerifyDependencies(identifiers); err != nil {
			return nil, err
		}

		// auto-complete URI from named connections
		if t.Command.Connection != "" {
			uri, err := c.Connections.GetURIByName(t.Command.Connection)

			if err != nil {
				return nil, err
			}

			t.Command.URI = uri
		}

		// use default connection if no URI is provided
		if t.Command.URI == "" {
			t.Command.URI, _ = c.Connections.GetURIByName("default")
		}

		// append task to final tasks
		if t.Command.Text != "" {
			finalTasks = append(finalTasks, t)
		}

		// parse queries from file and append new related tasks
		if t.Command.File != "" {
			parser, err := parser.NewParserBuilder(t.Command.Type).
				FromFile(t.Command.File).
				Build()

			if err != nil {
				return nil, err
			}

			for queryId, query := range parser.Parse() {
				finalTasks = append(finalTasks, task.Task{
					ID:      t.ID,
					QueryID: queryId,
					Name:    fmt.Sprintf("Query loaded from %s", t.Command.File),
					Command: task.Command{
						Text: query,
						Type: t.Command.Type,
						URI:  t.Command.URI,
					},
				})
			}
		}

		// append task to already knwown identifiers
		identifiers = append(identifiers, t.ID)
	}

	return finalTasks, nil
}
