package app

import (
	"time"
)

// isBankholiday determines if a date matches Cuvva's defintion of a bank holiday
func (a *App) isBankholiday(today time.Time) bool {
	tomorrow := today.Add(12 * time.Hour)
	yesterday := today.Add(-12 * time.Hour)

	if _, ok := a.bankholidays[yesterday.Format("2006-01-02")]; ok {
		return true
	}

	if _, ok := a.bankholidays[today.Format("2006-01-02")]; ok {
		return true
	}

	if _, ok := a.bankholidays[tomorrow.Format("2006-01-02")]; ok {
		return true
	}

	return false
}
