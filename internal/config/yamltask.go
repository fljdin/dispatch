package config

import (
	"github.com/fljdin/dispatch/internal/actions"
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
	Depends     []int      `yaml:"depends_on,omitempty"`
	Loader      YamlLoader `yaml:"loaded,omitempty"`
	Environment string     `yaml:"env,omitempty"`
	Variables   Variables  `yaml:"variables,omitempty"`
}

func (t YamlTask) Normalize(env Environments) (Task, error) {
	// make variables map if it's nil
	if t.Variables == nil {
		t.Variables = make(Variables)
	}

	// auto-complete environment variables
	if t.Environment != "" {
		env, err := env.ByName(t.Environment)

		if err != nil {
			return Task{}, err
		}

		// own variables take precedence over env variables
		t.Variables = t.Variables.Inherit(env.Variables)
	}

	// variables take precedence over default variables
	if env, err := env.ByName("default"); err == nil {
		t.Variables = t.Variables.Inherit(env.Variables)
	}

	// use shell as default type
	if t.Type == "" {
		t.Type = actions.Shell
	}

	var action actions.Actioner

	if !t.Loader.IsZero() {
		if t.Loader.Command != "" && t.Loader.From != "" {
			// auto-complete loader environment variables
			if t.Loader.Environment != "" {
				env, err := env.ByName(t.Loader.Environment)

				if err != nil {
					return Task{}, err
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

			action = actions.OutputLoader{
				Text: t.Loader.Command,
				From: t.Loader.From,
				Type: t.Type,
				Variables: actions.NestedVariables{
					Outer: t.Variables,
					Inner: t.Loader.Variables,
				},
			}
		}
	} else if t.File != "" {
		action = actions.FileLoader{
			File:      t.File,
			Type:      t.Type,
			Variables: t.Variables,
		}
	} else {
		action = actions.Command{
			Text:      t.Command,
			Type:      t.Type,
			Variables: t.Variables,
		}
	}

	return Task{
		Identifier: NewId(t.ID, 0),
		Name:       t.Name,
		Action:     action,
		Depends:    t.Depends,
	}, nil
}
