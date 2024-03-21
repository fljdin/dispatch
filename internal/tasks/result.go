package tasks

import (
	"time"
)

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

func (r Result) LoggerArgs() []any {
	return []any{
		"status", StatusTypes[r.Status],
		"elapsed", r.Elapsed.Round(time.Millisecond),
		"task", r.ID,
		"command", r.SubID,
		"worker", r.WorkerID,
	}
}
