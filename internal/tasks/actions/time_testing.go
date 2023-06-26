//go:build testing
// +build testing

package actions

import "time"

func Time() time.Time {
	return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
}
