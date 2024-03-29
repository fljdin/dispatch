package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigBuilder struct {
	config Config
	err    error
}

func NewBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (cb *ConfigBuilder) WithProcesses(value int) *ConfigBuilder {
	if value < 1 {
		return cb
	}

	cb.config.Processes = value
	return cb
}

func (cb *ConfigBuilder) WithTask(task YamlTask) *ConfigBuilder {
	cb.config.DeclaredTasks = append(cb.config.DeclaredTasks, task)
	return cb
}

func (cb *ConfigBuilder) WithLogfile(filename string) *ConfigBuilder {
	if len(filename) == 0 {
		return cb
	}

	cb.config.Logfile = filename
	return cb
}

func (cb *ConfigBuilder) WithYAML(yamlString string) *ConfigBuilder {
	if cb.err != nil {
		return cb
	}
	var config Config
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		cb.err = fmt.Errorf("error parsing yaml config: %w", err)
	}

	cb.config = config
	return cb
}

func (cb *ConfigBuilder) FromYAML(yamlFilename string) *ConfigBuilder {
	if cb.err != nil {
		return cb
	}
	if len(yamlFilename) == 0 {
		return cb
	}

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

	cb.config.ConfigureProcesses()

	return cb.config, cb.err
}
