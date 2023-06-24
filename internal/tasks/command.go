package tasks

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"golang.org/x/exp/slices"
)

var CommandTypes = []string{"sh", "psql"}
var testing bool = os.Getenv("GOTEST") != ""

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

func (c Command) VerifyType() error {
	if c.Type == "" {
		return nil
	}
	if !slices.Contains(CommandTypes, c.Type) {
		return fmt.Errorf("%s is not supported", c.Type)
	}
	return nil
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

func (Command) Time() time.Time {
	if testing {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Now()
}
