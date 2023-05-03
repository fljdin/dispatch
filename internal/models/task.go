package models

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"golang.org/x/exp/slices"
)

const (
	Waiting int = iota
	Succeeded
	Failed
	Interrupted
)

var TaskTypes = []string{"sh", "psql"}

type Task struct {
	ID         int    `yaml:"id"`
	Type       string `yaml:"type,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Command    string `yaml:"command"`
	File       string `yaml:"file"`
	URI        string `yaml:"uri,omitempty"`
	Connection string `yaml:"connection,omitempty"`
	Depends    []int  `yaml:"depends_on,omitempty"`
	Output     string `yaml:"output,omitempty"`
	QueryID    int
}

func (t Task) VerifyRequired() error {
	if t.ID == 0 {
		return fmt.Errorf("id is required")
	}

	if t.Command == "" && t.File == "" {
		return fmt.Errorf("command is required")
	}

	return nil
}

func (t Task) VerifyType() error {
	if "" == t.Type {
		return nil
	}
	if !slices.Contains(TaskTypes, t.Type) {
		return fmt.Errorf("%s is an invalid task type", t.Type)
	}
	return nil
}

func (t Task) VerifyDependencies(identifiers []int) error {
	verified := true

	for _, d := range t.Depends {
		verified = verified && slices.Contains(identifiers, d)
	}

	if !verified {
		return fmt.Errorf("task %d depends on unknown task %d", t.ID, t.Depends)
	}

	return nil
}

func (t Task) writeOutput(output []byte) {
	if len(t.Output) > 0 {
		_ = os.WriteFile(t.Output, output, 0644)
	}
}

func (t Task) Run() TaskResult {
	var cmd *exec.Cmd

	startTime := time.Now()

	switch t.Type {
	case "psql":
		cmd = exec.Command("psql", "-d", t.URI, "-c", t.Command)
	default:
		cmd = exec.Command("sh", "-c", t.Command)
	}

	output, err := cmd.CombinedOutput()
	endTime := time.Now()

	t.writeOutput(output)

	tr := TaskResult{
		ID:        t.ID,
		QueryID:   t.QueryID,
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
