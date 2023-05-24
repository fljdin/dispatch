package config

import (
	"fmt"
	"runtime"

	"github.com/fljdin/dispatch/internal/models"
	"github.com/fljdin/dispatch/internal/parser"
)

var ConfigWorkersDefault int = 2

type Config struct {
	Tasks             []models.Task      `yaml:"tasks"`
	MaxWorkers        int                `yaml:"workers"`
	Logfile           string             `yaml:"logfile"`
	Connections       models.Connections `yaml:"connections"`
	DefaultConnection models.Connection
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
		c.Connections = append(c.Connections, models.Connection{
			Name: "default",
			URI:  c.DefaultConnection.CombinedURI(),
		})
	}
}

func (c *Config) FinalizeTasks() ([]models.Task, error) {
	var finalTasks []models.Task
	var identifiers []int

	for _, t := range c.Tasks {
		if err := t.VerifyRequired(); err != nil {
			return nil, err
		}

		if err := t.VerifyType(); err != nil {
			return nil, err
		}

		if err := t.VerifyDependencies(identifiers); err != nil {
			return nil, err
		}

		// auto-complete URI from named connections
		if t.URI == "" && t.Connection != "" {
			uri, err := c.Connections.GetURIByName(t.Connection)

			if err != nil {
				return nil, err
			}

			t.URI = uri
		}

		// use default connection if no URI is provided
		if t.URI == "" {
			t.URI, _ = c.Connections.GetURIByName("default")
		}

		// append task to final tasks
		if t.Command != "" {
			finalTasks = append(finalTasks, t)
		}

		// parse queries from file and append new related tasks
		if t.Command == "" && t.File != "" {
			parser, err := parser.NewParserBuilder(t.Type).
				FromFile(t.File).
				Build()

			if err != nil {
				return nil, err
			}

			for queryId, query := range parser.Parse() {
				finalTasks = append(finalTasks, models.Task{
					ID:      t.ID,
					QueryID: queryId,
					Type:    t.Type,
					Name:    fmt.Sprintf("Query loaded from %s", t.File),
					Command: query,
					URI:     t.URI,
					Output:  t.Output,
				})
			}
		}

		// append task to already knwown identifiers
		identifiers = append(identifiers, t.ID)
	}

	return finalTasks, nil
}
