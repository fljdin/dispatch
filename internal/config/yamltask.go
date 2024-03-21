package config

import (
	"github.com/fljdin/dispatch/internal/tasks"
)

type YamlLoader struct {
	From        string    `yaml:"from"`
	Command     string    `yaml:"command"`
	Environment string    `yaml:"env,omitempty"`
	Variables   Variables `yaml:"variables,omitempty"`
}

func (y YamlLoader) IsZero() bool {
	return y.From == "" && y.Command == "" &&
		y.Environment == "" && y.Variables == nil
}

type YamlTask struct {
	ID          int        `yaml:"id"`
	Type        string     `yaml:"type,omitempty"`
	Name        string     `yaml:"name,omitempty"`
	Command     string     `yaml:"command"`
	File        string     `yaml:"file"`
	URI         string     `yaml:"uri,omitempty"`
	Connection  string     `yaml:"connection,omitempty"`
	Depends     []int      `yaml:"depends_on,omitempty"`
	Loader      YamlLoader `yaml:"loaded,omitempty"`
	Environment string     `yaml:"env,omitempty"`
	Variables   Variables  `yaml:"variables,omitempty"`
}

func (t YamlTask) Normalize(cnx Connections, env Environments) (tasks.Task, error) {

	// auto-complete URI from named connections
	if t.Connection != "" {
		uri, err := cnx.GetURIByName(t.Connection)

		if err != nil {
			return tasks.Task{}, err
		}

		t.URI = uri
	}

	// make variables map if it's nil
	if t.Variables == nil {
		t.Variables = make(Variables)
	}

	// auto-complete environment variables
	if t.Environment != "" {
		env, err := env.ByName(t.Environment)

		if err != nil {
			return tasks.Task{}, err
		}

		// own variables take precedence over env variables
		t.Variables = t.Variables.Inherit(env.Variables)
	}

	// variables take precedence over default variables
	if env, err := env.ByName("default"); err == nil {
		t.Variables = t.Variables.Inherit(env.Variables)
	}

	// use default connection if no URI is provided
	if t.URI == "" {
		t.URI, _ = cnx.GetURIByName("default")
	}

	var action tasks.Action

	if !t.Loader.IsZero() {
		if t.Loader.Command != "" && t.Loader.From != "" {
			// auto-complete loader environment variables
			if t.Loader.Environment != "" {
				env, err := env.ByName(t.Loader.Environment)

				if err != nil {
					return tasks.Task{}, err
				}

				// make variables map if it's nil
				if t.Loader.Variables == nil {
					t.Loader.Variables = make(Variables)
				}

				// own variables take precedence over env variables
				t.Loader.Variables = t.Loader.Variables.Inherit(env.Variables)
			}

			// inherit variables from task
			t.Loader.Variables = t.Loader.Variables.Inherit(t.Variables)

			action = tasks.OutputLoader{
				Text: t.Loader.Command,
				From: t.Loader.From,
				Type: t.Type,
				URI:  t.URI,
				Variables: tasks.NestedVariables{
					Outer: t.Variables,
					Inner: t.Loader.Variables,
				},
			}
		}
	} else if t.File != "" {
		action = tasks.FileLoader{
			File:      t.File,
			Type:      t.Type,
			URI:       t.URI,
			Variables: t.Variables,
		}
	} else {
		action = tasks.Command{
			Text:      t.Command,
			Type:      t.Type,
			URI:       t.URI,
			Variables: t.Variables,
		}
	}

	return tasks.Task{
		ID:      t.ID,
		Name:    t.Name,
		Action:  action,
		Depends: t.Depends,
	}, nil
}
