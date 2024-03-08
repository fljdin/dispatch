package tasks

import (
	"fmt"
	"os"

	"github.com/fljdin/fragment/languages"
	"golang.org/x/exp/slices"
)

type FileLoader struct {
	File string
	Type string
	URI  string
}

func (l FileLoader) load(input string) []string {
	switch l.Type {
	case "psql":
		return languages.PgSQL.Split(input)
	default:
		return languages.Shell.Split(input)
	}
}

func (l FileLoader) Validate() error {

	if !slices.Contains(CommandTypes, l.Type) {
		return fmt.Errorf("%s is not supported", l.Type)
	}

	if l.File != "" && l.Type == "" {
		return fmt.Errorf("type is required with a file")
	}

	if l.File == "" {
		return fmt.Errorf("file is required")
	}

	return nil
}

func (l FileLoader) Run() (Report, []Action) {
	startTime := Time()
	data, err := os.ReadFile(l.File)

	if err != nil {
		endTime := Time()

		result := Report{
			StartTime: startTime,
			EndTime:   endTime,
			Elapsed:   endTime.Sub(startTime),
			Status:    Failed,
			Error:     err.Error(),
		}
		return result, nil
	}

	var commands []Action
	for _, command := range l.load(string(data)) {
		commands = append(commands, Command{
			Text: command,
			Type: l.Type,
			URI:  l.URI,
		})
	}

	endTime := Time()

	result := Report{
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    Ready,
		Output:    fmt.Sprintf("%d loaded from %s", len(commands), l.File),
	}

	return result, commands
}
