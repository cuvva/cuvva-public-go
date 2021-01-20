package app

import (
	"context"
	"fmt"
)

func (a *App) GetBankHolidays(ctx context.Context, year int) error {
	holidays, err := a.govUK.GetBankHolidays(ctx)
	if err != nil {
		return err
	}

	m := map[string]struct{}{}

	for _, e := range holidays.EnglandAndWales.Events {
		m[e.Date] = struct{}{}
	}

	addExtras(m, year)

	a.bankholidays = m

	return nil
}

func addExtras(m map[string]struct{}, y int) {
	m[fmt.Sprintf("%d-12-26", y)] = struct{}{}
	m[fmt.Sprintf("%d-12-27", y)] = struct{}{}
	m[fmt.Sprintf("%d-12-28", y)] = struct{}{}
	m[fmt.Sprintf("%d-12-29", y)] = struct{}{}
	m[fmt.Sprintf("%d-12-30", y)] = struct{}{}
	m[fmt.Sprintf("%d-12-31", y)] = struct{}{}
}
