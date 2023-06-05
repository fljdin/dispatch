package models

import (
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/exp/slices"
)

var CommandTypes = []string{"sh", "psql"}

type Command struct {
	Text       string
	File       string
	Type       string
	URI        string
	Connection string
}

func (c Command) getExecCommand() *exec.Cmd {
	switch c.Type {
	case "psql":
		cmd := exec.Command("psql", "-d", c.URI)

		// use standard input to handle \g meta-commands
		textPipe, _ := cmd.StdinPipe()
		fmt.Fprintf(textPipe, c.Text)
		defer textPipe.Close()

		return cmd
	default:
		return exec.Command("sh", "-c", c.Text)
	}
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

func (c Command) Run() TaskResult {
	cmd := c.getExecCommand()

	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	endTime := time.Now()

	result := TaskResult{
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    Succeeded,
		Output:    string(output),
	}

	if err != nil {
		result.Status = Failed
		result.Error = err.Error()
	}

	return result
}
