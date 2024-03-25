package tasks

import (
	"time"

	"golang.org/x/exp/slices"
)

var CommandTypes = []string{"", "sh", "psql"}
var StatusTypes = []string{"waiting", "interrupted", "failed", "ready", "succeeded"}

func IsSucceeded(status int) bool {
	s := []int{Ready, Succeeded}
	return slices.Contains(s, status)
}

func IsFailed(status int) bool {
	s := []int{Interrupted, Failed}
	return slices.Contains(s, status)
}

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
	String() string
}
