package actions

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/parser"
	"golang.org/x/exp/slices"
)

type OutputLoader struct {
	Text string
	From string
	Type string
	URI  string
}

func (c OutputLoader) Validate() error {

	if c.Text == "" {
		return fmt.Errorf("command is required")
	}

	if !slices.Contains(CommandTypes, c.Type) {
		return fmt.Errorf("%s is not supported", c.Type)
	}

	if !slices.Contains(CommandTypes, c.From) {
		return fmt.Errorf("%s is not supported", c.From)
	}

	return nil
}

func (c OutputLoader) Run() (Report, []Action) {
	if c.From == "psql" {
		c.Text = fmt.Sprintf("%s \\g (format=unaligned tuples_only)", c.Text)
	}

	cmd := Command{
		Text: c.Text,
		Type: c.From,
		URI:  c.URI,
	}

	err := cmd.Validate()

	if err != nil {
		return Report{Status: KO, Error: err.Error()}, nil
	}

	result, _ := cmd.Run()

	if result.Status == KO {
		return result, nil
	}

	parser, err := parser.NewBuilder(c.Type).
		WithContent(result.Output).
		Build()

	if err != nil {
		result.Status = KO
		result.Error = err.Error()
		return result, nil
	}

	var commands []Action
	for _, command := range parser.Parse() {
		commands = append(commands, Command{
			Text: command,
			Type: c.Type,
			URI:  c.URI,
		})
	}

	return result, commands
}
