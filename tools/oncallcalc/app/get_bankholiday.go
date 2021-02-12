package app

import (
	"context"
)

func (a *App) GetBankHolidays(ctx context.Context) error {
	holidays, err := a.govUK.GetBankHolidays(ctx)
	if err != nil {
		return err
	}

	m := map[string]struct{}{}

	for _, e := range holidays.EnglandAndWales.Events {
		m[e.Date] = struct{}{}
	}

	a.bankholidays = m

	return nil
}
