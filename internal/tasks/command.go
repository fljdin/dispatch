package tasks

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/exp/slices"
)

type Command struct {
	Text      string
	Type      string
	URI       string
	Variables map[string]string
}

func (c Command) Validate() error {

	if !slices.Contains(CommandTypes, c.Type) {
		return fmt.Errorf("%s is not supported", c.Type)
	}

	if c.Text == "" {
		return fmt.Errorf("command is required")
	}

	return nil
}

func (c Command) getExecCommand() *exec.Cmd {
	var cmd *exec.Cmd

	switch c.Type {
	case "psql":
		// ON_ERROR_STOP is used to retrieve the correct exit code
		cmd = exec.Command("psql", "-v", "ON_ERROR_STOP=1", "-d", c.URI)

		// use input pipe to handle \g meta-commands
		textPipe, _ := cmd.StdinPipe()
		go func() {
			defer textPipe.Close()
			io.WriteString(textPipe, c.Text)
		}()
	default:
		cmd = exec.Command("sh", "-c", c.Text)
	}

	// set environment variables
	cmd.Env = os.Environ()
	for key, value := range c.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	return cmd
}

func (c Command) Run() (Report, []Action) {
	startTime := Time()

	var stdout, stderr bytes.Buffer
	cmd := c.getExecCommand()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	endTime := Time()

	result := Report{
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

	return result, nil
}
