package main

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Tasks       []Task      `yaml:"tasks"`
	MaxWorkers  int         `yaml:"workers"`
	Connections Connections `yaml:"connections"`
}

func (c *Config) ConfigureWorkers() {
	if c.MaxWorkers < 1 {
		c.MaxWorkers = ConfigWorkersDefault
	}

	if c.MaxWorkers > runtime.NumCPU() {
		c.MaxWorkers = runtime.NumCPU()
	}
}

func (c *Config) FinalizeTasks() ([]Task, error) {
	var finalTasks []Task

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
			parser, err := NewParserBuilder(t.Type).
				FromFile(t.File).
				Build()

			if err != nil {
				return nil, err
			}

			for _, query := range parser.Parse() {
				finalTasks = append(finalTasks, Task{
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

type ConfigBuilder struct {
	config Config
	err    error
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (cb *ConfigBuilder) WithMaxWorkers(value int) *ConfigBuilder {
	cb.config.MaxWorkers = value
	return cb
}

func (cb *ConfigBuilder) WithTask(task Task) *ConfigBuilder {
	cb.config.Tasks = append(cb.config.Tasks, task)
	return cb
}

func (cb *ConfigBuilder) WithYAML(yamlString string) *ConfigBuilder {
	var config Config
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		cb.err = fmt.Errorf("error parsing yaml config: %w", err)
	}

	cb.config = config
	return cb
}

func (cb *ConfigBuilder) FromYAML(yamlFilename string) *ConfigBuilder {
	data, err := os.ReadFile(yamlFilename)
	if err != nil {
		cb.err = fmt.Errorf("error reading yaml file: %w", err)
	}

	cb.config = cb.WithYAML(string(data)).config
	return cb
}

func (cb *ConfigBuilder) Build() (Config, error) {
	// if an error has already occurred, stop here
	if cb.err != nil {
		return cb.config, cb.err
	}

	cb.config.ConfigureWorkers()
	cb.config.Tasks, cb.err = cb.config.FinalizeTasks()

	return cb.config, cb.err
}
