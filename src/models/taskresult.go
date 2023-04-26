package models

import "time"

type TaskResult struct {
	ID        int
	Task      *Task
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Status    int
	Output    string
	Error     string
}
