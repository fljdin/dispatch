package tasks

import (
	"time"
)

var CommandTypes = []string{"", "sh", "psql"}

const (
	Waiting int = iota
	Interrupted
	Failed
	Ready
	Succeeded
)

type Report struct {
	Status    int
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Output    string
	Error     string
}

type Action interface {
	Validate() error
	Run() (Report, []Action)
}
