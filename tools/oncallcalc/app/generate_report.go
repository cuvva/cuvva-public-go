package app

import (
	"fmt"
	"time"

	"github.com/cuvva/cuvva-public-go/tools/oncallcalc"
)

// GenerateRota generates a rota based on a schedule and month
func (a *App) GenerateRota(scheduleID string, year int, month time.Month) (oncallcalc.Rota, interface{}, error) {
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}

	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	monthEnd := monthStart.AddDate(0, 1, 0)

	pdshifts, badShift, err := a.ListPagerDutyShifts(scheduleID, monthStart, monthEnd)
	if err != nil {
		return nil, badShift, err
	}

	rota := makeRota(pdshifts)
	shifts := a.ConvertShifts(pdshifts)

	for _, shift := range shifts {
		if shift.Date.Month() != month {
			continue
		}

		v, ok := rota[shift.Email]
		if !ok {
			return nil, rota, fmt.Errorf("%s is missing from the rota", shift.Email)
		}

		if isWeekend(shift.Date) {
			v.Weekends += 0.5
		} else {
			v.Weekdays += 0.5
		}
	}

	return rota, nil, nil
}

// makeRota builds up a Rota struct based on names in the PagerDuty shifts
func makeRota(in []*PDShift) oncallcalc.Rota {
	r := make(oncallcalc.Rota)

	for _, v := range in {
		r[v.Email] = &oncallcalc.Stat{}
	}

	return r
}
