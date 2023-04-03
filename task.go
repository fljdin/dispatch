package main

import (
	"os/exec"
	"time"
)

const (
	Failed int = iota
	Succeeded
)

type TaskResult struct {
	ID        int
	StartTime time.Time
	EndTime   time.Time
	Status    int
	Message   string
}

type Task struct {
	ID      int    `yaml:"id"`
	Command string `yaml:"command"`
}

func (t Task) Run(results chan<- TaskResult) {
	startTime := time.Now()

	cmd := exec.Command("sh", "-c", t.Command)

	output, err := cmd.CombinedOutput()
	endTime := time.Now()

	tr := TaskResult{
		ID:        t.ID,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    Succeeded,
		Message:   string(output),
	}

	if err != nil {
		tr.Message = err.Error()
	}

	results <- tr
}
