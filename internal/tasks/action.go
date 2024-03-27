package tasks

import (
	"time"

	"github.com/fljdin/dispatch/internal/status"
)

var CommandTypes = []string{"", "sh", "psql"}

type Report struct {
	Status    status.Status
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Output    string
	Error     string
}

type Action interface {
	Validate() error
	Run() (Report, []Action)
	String() string
}
