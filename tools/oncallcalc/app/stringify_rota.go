package app

import (
	"fmt"

	"github.com/cuvva/cuvva-public-go/tools/oncallcalc"
)

func (a *App) StringifyRota(rota oncallcalc.Rota) string {
	var out string
	var total float64

	for e, s := range rota {
		out += fmt.Sprintf("%s\n", e)
		out += fmt.Sprintf("\tDay count (wd/we/bh): (%.1f, %.1f, %.1f)\n", s.Weekdays, s.Weekends, s.Bankholidays)

		bankholidayMoney := s.Bankholidays * oncallcalc.WeekendPayout
		weekdayMoney := s.Weekdays * oncallcalc.WeekdayPayout
		weekendMoney := s.Weekends * oncallcalc.WeekendPayout
		total += weekendMoney + weekdayMoney + bankholidayMoney

		out += fmt.Sprintf("\tPayout (wd/we/bh): (%.2f, %.2f, %.2f)\n", weekdayMoney, weekendMoney, bankholidayMoney)

		out += fmt.Sprintf("\tTotal: £%.2f\n\n", weekdayMoney+weekendMoney+bankholidayMoney)
	}

	out += fmt.Sprintf("Total: £%.2f\n", total)

	return out
}
