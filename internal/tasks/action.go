package tasks

import (
	"time"

	"github.com/fljdin/dispatch/internal/status"
)

var (
	Shell        = "sh"
	PgSQL        = "psql"
	CommandTypes = []string{Shell, PgSQL}
)

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
	Command() string
	String() string
}
