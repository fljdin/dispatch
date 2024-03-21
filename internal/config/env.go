package config

import (
	"fmt"
)

type Environment struct {
	Name      string            `yaml:"name"`
	Variables map[string]string `yaml:"variables"`
}

type Environments []Environment

func (e Environments) ByName(name string) (Environment, error) {
	for _, env := range e {
		if env.Name == name {
			return env, nil
		}
	}

	return Environment{}, fmt.Errorf("environment not found for name %s", name)
}
