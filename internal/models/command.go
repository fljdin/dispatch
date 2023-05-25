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
	var cmd *exec.Cmd

	startTime := time.Now()

	switch c.Type {
	case "psql":
		cmd = exec.Command("psql", "-d", c.URI, "-c", c.Text)
	default:
		cmd = exec.Command("sh", "-c", c.Text)
	}

	output, err := cmd.CombinedOutput()
	endTime := time.Now()

	tr := TaskResult{
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    Succeeded,
		Output:    string(output),
	}

	if err != nil {
		tr.Status = Failed
		tr.Error = err.Error()
	}

	return tr
}
