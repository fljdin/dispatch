package tasks

import (
	"time"
)

type Result struct {
	ID        int
	Name      string
	Action    string
	SubID     int
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Status    int
	Output    string
	Error     string
}
