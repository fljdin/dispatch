//go:build !testing
// +build !testing

package tasks

import "time"

func Time() time.Time {
	return time.Now()
}
