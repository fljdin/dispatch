package config

import (
	"fmt"
	"os"

	"github.com/fljdin/dispatch/internal/tasks"
	"gopkg.in/yaml.v2"
)

const ProcessesDefault int = 1

type Config struct {
	DeclaredTasks []YamlTask   `yaml:"tasks"`
	Environments  Environments `yaml:"environments"`
}

func New(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error reading yaml file: %w", err)
	}

	return NewFromRaw(string(data))
}

func NewFromRaw(raw string) (Config, error) {
	var config Config
	err := yaml.Unmarshal([]byte(raw), &config)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing yaml config: %w", err)
	}

	return config, nil
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
