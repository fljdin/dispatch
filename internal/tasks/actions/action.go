package actions

import (
	"os"
	"time"
)

var CommandTypes = []string{"", "sh", "psql"}
var testing bool = os.Getenv("GOTEST") != ""

const (
	OK int = iota + 2
	KO
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

func Time() time.Time {
	if testing {
		return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Now()
}
