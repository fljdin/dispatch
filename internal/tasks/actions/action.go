package actions

import (
	"time"
)

var CommandTypes = []string{"", "sh", "psql"}

const (
	KO int = iota + 2
	Loaded
	OK
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
