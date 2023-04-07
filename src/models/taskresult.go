package models

import "time"

type TaskResult struct {
	ID        int
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Status    int
	Output    string
	Error     string
}
