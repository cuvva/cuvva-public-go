package oncallcalc

import (
	"time"
)

// WeekendPayout is the amount of money in GBP to pay for a weekend shift
const WeekendPayout = 100

// WeekdayPayout is the amount of money in GBP to pay for a weekday shift
const WeekdayPayout = 50

// Rota is the main data object storing everyones on-call shift counts
type Rota map[string]*Stat

// Stat is a gives out the day stats for 1 person
type Stat struct {
	Weekdays float64
	Weekends float64
}

// Shift is a Cuvva specific representation of a 12 hour on-call shift
type Shift struct {
	Email string
	Date  time.Time
}
