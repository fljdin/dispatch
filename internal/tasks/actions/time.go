//go:build !testing
// +build !testing

package actions

import "time"

func Time() time.Time {
	return time.Now()
}
