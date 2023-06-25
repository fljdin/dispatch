package config

import (
	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/fljdin/dispatch/internal/tasks/actions"
)

type YamlLoader struct {
	From    string `yaml:"from"`
	Command string `yaml:"command"`
	File    string `yaml:"file"`
}

type YamlTask struct {
	ID         int        `yaml:"id"`
	Type       string     `yaml:"type,omitempty"`
	Name       string     `yaml:"name,omitempty"`
	Command    string     `yaml:"command"`
	URI        string     `yaml:"uri,omitempty"`
	Connection string     `yaml:"connection,omitempty"`
	Depends    []int      `yaml:"depends_on,omitempty"`
	Generated  YamlLoader `yaml:"generated,omitempty"`
}

func (t YamlTask) Normalize(cnx tasks.Connections) (tasks.Task, error) {

	// auto-complete URI from named connections
	if t.Connection != "" {
		uri, err := cnx.GetURIByName(t.Connection)

		if err != nil {
			return tasks.Task{}, err
		}

		t.URI = uri
	}

	// use default connection if no URI is provided
	if t.URI == "" {
		t.URI, _ = cnx.GetURIByName("default")
	}

	var action actions.Action

	if t.Generated != (YamlLoader{}) {
		if t.Generated.File != "" {
			action = actions.FileLoader{
				File: t.Generated.File,
				Type: t.Type,
				URI:  t.URI,
			}
		} else if t.Generated.Command != "" && t.Generated.From != "" {
			action = actions.OutputLoader{
				Text: t.Generated.Command,
				From: t.Generated.From,
				Type: t.Type,
				URI:  t.URI,
			}
		}
	} else {
		action = actions.Command{
			Text: t.Command,
			Type: t.Type,
			URI:  t.URI,
		}
	}

	return tasks.Task{
		ID:      t.ID,
		Name:    t.Name,
		Action:  action,
		Depends: t.Depends,
	}, nil
}
