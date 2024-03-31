package actions

import (
	"fmt"
	"os"

	"github.com/fljdin/dispatch/internal/helper"
	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/fragment/languages"
	"golang.org/x/exp/slices"
)

type File struct {
	File      string
	Type      string
	Variables map[string]string
}

func (l File) load(input string) []string {
	switch l.Type {
	case PgSQL:
		return languages.PgSQL.Split(input)
	default:
		return languages.Shell.Split(input)
	}
}

func (l File) String() string {
	return fmt.Sprintf("execute %s", l.File)
}

func (l File) Command() string {
	return l.Type
}

func (l File) Validate() error {

	if !slices.Contains(CommandTypes, l.Type) {
		return fmt.Errorf("%s is not supported", l.Type)
	}

	if l.File == "" {
		return fmt.Errorf("file is required")
	}

	return nil
}

func (l File) Run() (Result, []Actioner) {
	startTime := helper.Now()
	data, err := os.ReadFile(l.File)

	if err != nil {
		endTime := helper.Now()

		result := Result{
			StartTime: startTime,
			EndTime:   endTime,
			Elapsed:   endTime.Sub(startTime),
			Status:    status.Failed,
			Error:     err.Error(),
		}
		return result, nil
	}

	var commands []Actioner
	for _, command := range l.load(string(data)) {
		commands = append(commands, Command{
			Text:      command,
			Type:      l.Type,
			Variables: l.Variables,
		})
	}

	endTime := helper.Now()

	result := Result{
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    status.Succeeded,
		Output:    fmt.Sprintf("%d loaded from %s", len(commands), l.File),
	}

	return result, commands
}
