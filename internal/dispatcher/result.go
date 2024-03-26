package dispatcher

import "github.com/fljdin/dispatch/internal/status"

type Result struct {
	ID     int
	SubID  int
	Status status.Status
}
