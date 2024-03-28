//go:build testing
// +build testing

package helper

import "time"

func Now() time.Time {
	return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
}
