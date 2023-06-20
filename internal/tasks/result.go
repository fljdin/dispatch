package tasks

import "time"

type Result struct {
	ID        int
	SubID     int
	WorkerID  int
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Status    int
	Output    string
	Error     string
}
