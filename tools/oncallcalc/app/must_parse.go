package app

import (
	"time"
)

// MustParse parses a date into a time.Time, specifically for Europe/London
func MustParse(format, in string) time.Time {
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}

	f, err := time.Parse(format, in)
	if err != nil {
		panic(err)
	}

	return f.In(loc)
}
