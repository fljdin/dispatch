package config

import (
	"fmt"
)

type Variables map[string]string

// Inherit merges two variables maps, giving precedence to the original map
func (v Variables) Inherit(other Variables) Variables {
	if v == nil {
		v = make(Variables)
	}

	for k, val := range other {
		if _, ok := v[k]; !ok {
			v[k] = val
		}
	}
	return v
}

type Environment struct {
	Name      string    `yaml:"name"`
	Variables Variables `yaml:"variables"`
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
