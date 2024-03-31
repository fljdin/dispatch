package actions

import (
	"time"

	"github.com/fljdin/dispatch/internal/status"
)

type Result struct {
	Status    status.Status
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
	Output    string
	Error     string
}
