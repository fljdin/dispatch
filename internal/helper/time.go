//go:build !testing
// +build !testing

package helper

import "time"

func Now() time.Time {
	return time.Now()
}
