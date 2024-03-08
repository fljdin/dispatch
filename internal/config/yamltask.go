package config

import (
	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/fljdin/dispatch/internal/tasks/actions"
)

type YamlLoader struct {
	From    string `yaml:"from"`
	Command string `yaml:"command"`
}

type YamlTask struct {
	ID         int        `yaml:"id"`
	Type       string     `yaml:"type,omitempty"`
	Name       string     `yaml:"name,omitempty"`
	Command    string     `yaml:"command"`
	File       string     `yaml:"file"`
	URI        string     `yaml:"uri,omitempty"`
	Connection string     `yaml:"connection,omitempty"`
	Depends    []int      `yaml:"depends_on,omitempty"`
	Loader     YamlLoader `yaml:"loaded,omitempty"`
}

func (t YamlTask) Normalize(cnx Connections) (tasks.Task, error) {

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

	if t.Loader != (YamlLoader{}) {
		if t.Loader.Command != "" && t.Loader.From != "" {
			action = actions.OutputLoader{
				Text: t.Loader.Command,
				From: t.Loader.From,
				Type: t.Type,
				URI:  t.URI,
			}
		}
	} else if t.File != "" {
		action = actions.FileLoader{
			File: t.File,
			Type: t.Type,
			URI:  t.URI,
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
