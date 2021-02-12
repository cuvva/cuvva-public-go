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

func (d Date) Time(loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

func (d Date) String() string {
	return d.Time(time.UTC).Format("2006-01-02")
}
