package actions

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/parser"
	"golang.org/x/exp/slices"
)

type FileLoader struct {
	File string
	Type string
	URI  string
}

func (c FileLoader) Validate() error {

	if !slices.Contains(CommandTypes, c.Type) {
		return fmt.Errorf("%s is not supported", c.Type)
	}

	if c.File != "" && c.Type == "" {
		return fmt.Errorf("type is required with a file")
	}

	if c.File == "" {
		return fmt.Errorf("file is required")
	}

	return nil
}

func (c FileLoader) Run() (Report, []Action) {
	startTime := Time()

	parser, err := parser.NewBuilder(c.Type).
		FromFile(c.File).
		Build()

	if err != nil {
		endTime := Time()

		result := Report{
			StartTime: startTime,
			EndTime:   endTime,
			Elapsed:   endTime.Sub(startTime),
			Status:    KO,
			Error:     err.Error(),
		}
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

	endTime := Time()

	result := Report{
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    OK,
		Output:    fmt.Sprintf("%d loaded from %s", len(commands), c.File),
	}

	return result, commands
}
