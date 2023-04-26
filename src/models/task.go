package models

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

const (
	Waiting int = iota
	Ready
	Failed
	Succeeded
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
	Depends    []int  `yaml:"depends,omitempty"`
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

func (t Task) Run(ctx context.Context) TaskResult {
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
