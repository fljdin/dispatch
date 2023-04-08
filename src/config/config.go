package config

import (
	"fmt"
	"runtime"

	"github.com/fljdin/dispatch/src/models"
	"github.com/fljdin/dispatch/src/parser"
)

var ConfigWorkersDefault int = 2

type Config struct {
	Tasks       []models.Task      `yaml:"tasks"`
	MaxWorkers  int                `yaml:"workers"`
	Connections models.Connections `yaml:"connections"`
}

func (c *Config) ConfigureWorkers() {
	if c.MaxWorkers < 1 {
		c.MaxWorkers = ConfigWorkersDefault
	}

	if c.MaxWorkers > runtime.NumCPU() {
		c.MaxWorkers = runtime.NumCPU()
	}
}

func (c *Config) FinalizeTasks() ([]models.Task, error) {
	var finalTasks []models.Task

	for _, t := range c.Tasks {
		if err := t.VerifyRequired(); err != nil {
			return nil, err
		}

		if err := t.VerifyType(); err != nil {
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

			for _, query := range parser.Parse() {
				finalTasks = append(finalTasks, models.Task{
					ID:      t.ID,
					Type:    t.Type,
					Name:    fmt.Sprintf("Query loaded from %s", t.File),
					Command: query,
					URI:     t.URI,
				})
			}
		}
	}

	return finalTasks, nil
}
