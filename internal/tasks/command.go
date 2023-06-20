package tasks

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/fljdin/dispatch/internal/parser"
	"golang.org/x/exp/slices"
)

var CommandTypes = []string{"sh", "psql"}
var testing bool = os.Getenv("GOTEST") != ""

type Command struct {
	Text       string
	File       string
	From       string
	Type       string
	URI        string
	Connection string
}

func (Command) Time() time.Time {
	if testing {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Now()
}

func (c Command) VerifyType() error {
	if c.Type == "" {
		return nil
	}
	if !slices.Contains(CommandTypes, c.Type) {
		return fmt.Errorf("%s is not supported", c.Type)
	}
	return nil
}

func (c Command) getExecCommand() *exec.Cmd {
	switch c.Type {
	case "psql":
		// ON_ERROR_STOP is used to retrieve the correct exit code
		cmd := exec.Command("psql", "-v", "ON_ERROR_STOP=1", "-d", c.URI)

		// use input pipe to handle \g meta-commands
		textPipe, _ := cmd.StdinPipe()
		go func() {
			defer textPipe.Close()
			io.WriteString(textPipe, c.Text)
		}()

		return cmd
	default:
		return exec.Command("sh", "-c", c.Text)
	}
}

func (c Command) Run() Result {
	startTime := c.Time()

	var stdout, stderr bytes.Buffer
	cmd := c.getExecCommand()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	endTime := c.Time()

	result := Result{
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    Succeeded,
		Output:    stdout.String(),
		Error:     stderr.String(),
	}

	if err != nil {
		result.Status = Failed
	}

	return result
}

func (c Command) Generate() (Result, []Command) {
	var commands []Command
	result := Result{Status: Succeeded}

	parser, err := c.parser()
	if err != nil {
		result.Status = Failed
		result.Error = err.Error()
		return result, commands
	}

	for _, command := range parser.Parse() {
		commands = append(commands, Command{
			Text: command,
			Type: c.From,
		})
	}

	return result, commands
}

func (c Command) parser() (parser.Parser, error) {

	if c.File != "" {
		return parser.NewBuilder(c.From).
			FromFile(c.File).
			Build()
	}

	if c.Text != "" {
		result := c.Run()
		if result.Status == Failed {
			return nil, fmt.Errorf("%s", result.Error)
		}

		return parser.NewBuilder(c.From).
			WithContent(result.Output).
			Build()
	}

	return nil, fmt.Errorf("invalid generator")
}
