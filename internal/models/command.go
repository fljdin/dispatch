package models

import (
	"os/exec"
	"time"
)

type Command struct {
	Type    string `yaml:"type"`
	Command string `yaml:"command"`
}

func (c Command) Run() TaskResult {
	var cmd *exec.Cmd
	var (
		ID      int = 0
		QueryID int = 0
	)

	startTime := time.Now()

	switch c.Type {
	default:
		cmd = exec.Command("sh", "-c", c.Command)
	}

	output, _ := cmd.CombinedOutput()
	endTime := time.Now()

	tr := TaskResult{
		ID:        ID,
		QueryID:   QueryID,
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    Succeeded,
		Output:    string(output),
	}

	return tr
}
