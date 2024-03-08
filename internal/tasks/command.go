package tasks

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"golang.org/x/exp/slices"
)

type Command struct {
	Text string
	Type string
	URI  string
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
