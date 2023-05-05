package config

import (
	"fmt"
	"os"

	"github.com/fljdin/dispatch/internal/models"
	"gopkg.in/yaml.v2"
)

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

func (cb *ConfigBuilder) WithTask(task models.Task) *ConfigBuilder {
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