package config

import (
	"runtime"

	"github.com/fljdin/dispatch/internal/tasks"
)

var ConfigWorkersDefault int = 2

type Config struct {
	DeclaredTasks      []YamlTask   `yaml:"tasks"`
	MaxWorkers         int          `yaml:"workers"`
	Logfile            string       `yaml:"logfile"`
	Environments       Environments `yaml:"environments"`
	DefaultEnvironment Environment
}

func (c *Config) ConfigureWorkers() {
	if c.MaxWorkers < 1 {
		c.MaxWorkers = ConfigWorkersDefault
	}

	if c.MaxWorkers > runtime.NumCPU() {
		c.MaxWorkers = runtime.NumCPU()
	}
}

func (c Config) Tasks() ([]tasks.Task, error) {
	var finalTasks []tasks.Task
	var identifiers []int

	for _, declared := range c.DeclaredTasks {
		task, err := declared.Normalize(c.Environments)
		if err != nil {
			return nil, err
		}

		if err := task.Validate(); err != nil {
			return nil, err
		}

		if err := task.ValidateDependencies(identifiers); err != nil {
			return nil, err
		}

		finalTasks = append(finalTasks, task)
		identifiers = append(identifiers, task.ID)
	}

	return finalTasks, nil
}
