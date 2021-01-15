package app

import (
	"time"

	"github.com/cuvva/cuvva-public-go/tools/oncallcalc"
	log "github.com/sirupsen/logrus"
)

// ConvertShifts converts a slice of PagerDuty shifts to Cuvva shifts
func (a *App) ConvertShifts(pdshifts []*PDShift) []*oncallcalc.Shift {
	var all []*oncallcalc.Shift

	for _, s := range pdshifts {
		if !(s.Start.Hour()%12 == 0 || s.End.Hour()%12 == 0) {
			log.Warnf("shift does not align to 12 hour shift: %+v", s)
			continue
		}

		for next := s.Start; next.Before(s.End); next = next.Add(12 * time.Hour) {
			all = append(all, &oncallcalc.Shift{
				Email: s.Email,
				Date:  next,
			})
		}
	}

	return all
}
