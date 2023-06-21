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

var CommandTypes = []string{"", "sh", "psql"}
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

func (c Command) Validate() error {

	if !slices.Contains(CommandTypes, c.Type) {
		return fmt.Errorf("%s is not supported", c.Type)
	}

	if c.Text == "" && c.File == "" {
		return fmt.Errorf("command or file are required")
	}

	if c.From != "" && !slices.Contains(CommandTypes, c.From) {
		return fmt.Errorf("%s is not supported", c.From)
	}

	return nil
}

func (c Command) getExecCommand() *exec.Cmd {
	cmdType := c.Type
	if c.From != "" {
		cmdType = c.From
	}

	switch cmdType {
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
	var result Result
	var cmds []string

	if err := c.Validate(); err != nil {
		result = Result{
			Status: Failed,
			Error:  err.Error(),
		}
		return result, nil
	}

	if c.File != "" {
		result, cmds = c.generateFromFile()
	} else if c.From != "" && c.Text != "" {
		result, cmds = c.generateFromOutput()
	}

	var commands []Command
	for _, command := range cmds {
		commands = append(commands, Command{
			Text: command,
			Type: c.Type,
			URI:  c.URI,
		})
	}

	return result, commands
}

func (c Command) generateFromOutput() (Result, []string) {
	result := c.Run()

	if result.Status == Failed {
		return result, nil
	}

	parser, err := parser.NewBuilder(c.From).
		WithContent(result.Output).
		Build()

	if err != nil {
		result.Status = Failed
		result.Error = err.Error()
		return result, nil
	}

	cmds := parser.Parse()
	return result, cmds
}

func (c Command) generateFromFile() (Result, []string) {
	startTime := time.Now()

	parser, err := parser.NewBuilder(c.Type).
		FromFile(c.File).
		Build()

	if err != nil {
		endTime := time.Now()

		result := Result{
			StartTime: startTime,
			EndTime:   endTime,
			Elapsed:   endTime.Sub(startTime),
			Status:    Failed,
			Error:     err.Error(),
		}
		return result, nil
	}

	cmds := parser.Parse()
	endTime := time.Now()

	result := Result{
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    Succeeded,
		Output:    fmt.Sprintf("%d loaded from %s", len(cmds), c.File),
	}

	return result, cmds
}
