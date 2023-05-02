package models

import (
	"fmt"
	"os/exec"
	"time"
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
	for _, tt := range TaskTypes {
		if t.Type == tt || t.Type == "" {
			return nil
		}
	}
	return fmt.Errorf("%s is an invalid task type", t.Type)
}

func (t Task) VerifyDependencies(identifiers []int) error {
	verified := true

	for _, d := range t.Depends {
		found := false

		for _, i := range identifiers {
			if d == i {
				found = true
				break
			}
		}

		verified = verified && found
	}

	if !verified {
		return fmt.Errorf("task %d depends on unknown task %d", t.ID, t.Depends)
	}

	return nil
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

	tr := TaskResult{
		ID:        t.ID,
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
