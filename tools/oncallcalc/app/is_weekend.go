package app

import (
	"time"
)

// isWeekend determines if a date matches Cuvva's defintion of a weekend
func isWeekend(t time.Time) bool {
	if t.Weekday() == time.Friday && t.Hour() >= 12 {
		return true
	}

	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return true
	}

	if t.Weekday() == time.Monday && t.Hour() < 12 {
		return true
	}

	return false
}
