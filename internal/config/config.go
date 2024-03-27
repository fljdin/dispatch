package config

import (
	"runtime"

	"github.com/fljdin/dispatch/internal/tasks"
)

const ProcessesDefault int = 2

type Config struct {
	DeclaredTasks      []YamlTask   `yaml:"tasks"`
	Environments       Environments `yaml:"environments"`
	Logfile            string       `yaml:"logfile"`
	Processes          int          `yaml:"procs"`
	Verbose            bool         `yaml:"verbose"`
	DefaultEnvironment Environment
}

func (c *Config) ConfigureProcesses() {
	if c.Processes < 1 {
		c.Processes = ProcessesDefault
	}

	if c.Processes > runtime.NumCPU() {
		c.Processes = runtime.NumCPU()
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
		identifiers = append(identifiers, task.Identifier.ID)
	}

	return finalTasks, nil
}
