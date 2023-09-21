package app

import (
	"time"
)

// isWeekend determines if a times day is a Saturday or Sunday
func isWeekend(t time.Time) bool {

	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return true
	}

	return false
}
