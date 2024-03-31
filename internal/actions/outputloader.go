package actions

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/fragment/languages"
	"golang.org/x/exp/slices"
)

type NestedVariables struct {
	Outer map[string]string
	Inner map[string]string
}

type OutputLoader struct {
	Text      string
	From      string
	Type      string
	Variables NestedVariables
}

func (l OutputLoader) load(input string) []string {
	switch l.Type {
	case "psql":
		return languages.PgSQL.Split(input)
	default:
		return languages.Shell.Split(input)
	}
}

func (l OutputLoader) String() string {
	return l.Text
}

func (l OutputLoader) Command() string {
	return l.From
}

func (l OutputLoader) Validate() error {

	if l.Text == "" {
		return fmt.Errorf("command is required")
	}

	if !slices.Contains(CommandTypes, l.Type) {
		return fmt.Errorf("%s is not supported", l.Type)
	}

	if !slices.Contains(CommandTypes, l.From) {
		return fmt.Errorf("%s is not supported", l.From)
	}

	return nil
}

func (l OutputLoader) Run() (Result, []Actioner) {
	if l.From == "psql" {
		l.Text = fmt.Sprintf("%s \\g (format=unaligned tuples_only)", l.Text)
	}

	// run command with inner variables
	cmd := Command{
		Text:      l.Text,
		Type:      l.From,
		Variables: l.Variables.Inner,
	}

	err := cmd.Validate()

	if err != nil {
		return Result{Status: status.Failed, Error: err.Error()}, nil
	}

	result, _ := cmd.Run()

	if result.Status == status.Failed {
		return result, nil
	}

	var commands []Actioner
	for _, command := range l.load(result.Output) {
		// pass outer variables to children
		commands = append(commands, Command{
			Text:      command,
			Type:      l.Type,
			Variables: l.Variables.Outer,
		})
	}

	result.Status = status.Succeeded
	return result, commands
}
