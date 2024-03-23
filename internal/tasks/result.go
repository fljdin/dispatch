package tasks

import (
	"fmt"
	"time"
)

type Result struct {
	ID        int
	Name      string
	Action    string
	SubID     int
	ProcID    int
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Status    int
	Output    string
	Error     string
}

func (r Result) Code() string {
	return fmt.Sprintf("[%d:%d]", r.ID, r.SubID)
}

func (r Result) LoggerArgs() []any {
	return []any{
		"status", StatusTypes[r.Status],
		"name", r.Name,
		"elapsed", r.Elapsed.Round(time.Millisecond),
		"proc", r.ProcID,
	}
}
