package tasks

import (
	"fmt"

	"github.com/fljdin/fragment/languages"
	"golang.org/x/exp/slices"
)

type OutputLoader struct {
	Text string
	From string
	Type string
	URI  string
}

func (l OutputLoader) load(input string) []string {
	switch l.Type {
	case "psql":
		return languages.PgSQL.Split(input)
	default:
		return languages.Shell.Split(input)
	}
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

func (l OutputLoader) Run() (Report, []Action) {
	if l.From == "psql" {
		l.Text = fmt.Sprintf("%s \\g (format=unaligned tuples_only)", l.Text)
	}

	cmd := Command{
		Text: l.Text,
		Type: l.From,
		URI:  l.URI,
	}

	err := cmd.Validate()

	if err != nil {
		return Report{Status: Failed, Error: err.Error()}, nil
	}

	result, _ := cmd.Run()

	if result.Status == Failed {
		return result, nil
	}

	var commands []Action
	for _, command := range l.load(result.Output) {
		commands = append(commands, Command{
			Text: command,
			Type: l.Type,
			URI:  l.URI,
		})
	}

	result.Status = Ready
	return result, commands
}
