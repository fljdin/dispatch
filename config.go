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
	if cb.config.MaxWorkers < 1 {
		cb.config.MaxWorkers = ConfigWorkersDefault
	}

	if cb.config.MaxWorkers > runtime.NumCPU() {
		cb.config.MaxWorkers = runtime.NumCPU()
	}

	var finalTasks []Task

	for i, t := range cb.config.Tasks {
		if err := t.VerifyRequired(); err != nil {
			cb.err = err
			break
		}

		if err := t.VerifyType(); err != nil {
			cb.err = err
			break
		}

		// auto-complete URI from named connections
		if t.URI == "" && t.Connection != "" {
			if uri, err := cb.config.Connections.GetURIByName(t.Connection); err != nil {
				cb.err = err
			} else {
				t.URI = uri
				cb.config.Tasks[i] = t
			}
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
				cb.err = err
			} else {
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
	}

	cb.config.Tasks = finalTasks

	return cb.config, cb.err
}
