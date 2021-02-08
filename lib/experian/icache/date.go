package icache

import (
	"time"
)

type Date struct {
	Year  int        `xml:"CCYY"`
	Month time.Month `xml:"MM"`
	Day   int        `xml:"DD"`
}

func NewDate(t time.Time) Date {
	y, m, d := t.Date()

	return Date{y, m, d}
}
