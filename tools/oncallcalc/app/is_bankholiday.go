package app

import (
	"time"
)

// isBankholiday determines if a date is within the list of bankholidays
func (a *App) isBankholiday(today time.Time) bool {

	if _, ok := a.bankholidays[today.Format("2006-01-02")]; ok {
		return true
	}

	return false
}
