package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

type ConfigBuilder struct {
	config Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (cb *ConfigBuilder) WithTask(task Task) *ConfigBuilder {
	cb.config.Tasks = append(cb.config.Tasks, task)
	return cb
}

func (cb *ConfigBuilder) WithYAML(yamlString string) *ConfigBuilder {
	var config Config
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing YAML config: ", err)
		os.Exit(1)
	}

	cb.config = config
	return cb
}

func (cb *ConfigBuilder) FromYAML(yamlFilename string) (*ConfigBuilder, error) {
	data, err := ioutil.ReadFile(yamlFilename)
	if err != nil {
		return nil, fmt.Errorf("Error reading YAML file: %s", err)
	}

	return cb.WithYAML(string(data)), nil
}

func (cb *ConfigBuilder) Build() Config {
	return cb.config
}
