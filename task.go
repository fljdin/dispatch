package main

import (
	"context"
	"os/exec"
	"time"
)

const (
	Failed int = iota
	Succeeded
)

type Connection struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
}

type Task struct {
	ID         int    `yaml:"id"`
	Type       string `yaml:"type,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Command    string `yaml:"command"`
	Connection string `yaml:"connection,omitempty"`
}

type TaskResult struct {
	ID        int
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Status    int
	Message   string
}

func (t Task) Run(ctx context.Context) TaskResult {
	var cmd *exec.Cmd

	startTime := time.Now()

	switch t.Type {
	case "psql":
		cmd = exec.Command("psql", "-d", t.Connection, "-c", t.Command)
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
		Message:   string(output),
	}

	if err != nil {
		tr.Status = Failed
		tr.Message = err.Error()
	}

	return tr
}
