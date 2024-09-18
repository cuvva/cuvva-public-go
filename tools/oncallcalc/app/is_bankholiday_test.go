package app

import (
	"testing"
	"time"
)

var testData = []struct {
	Name        string
	InDate      time.Time
	BankHoliday bool
}{
	{
		Name:        "Not a bankholiday, nor near one",
		InDate:      mustParse("2020-03-08T00:00:00Z"),
		BankHoliday: false,
	},
	{
		Name:        "Certainly a bank holiday",
		InDate:      mustParse("2021-01-01T00:00:00Z"),
		BankHoliday: true,
	},
	{
		Name:        "Mayday: shift before",
		InDate:      mustParse("2021-05-02T00:00:00Z"),
		BankHoliday: false,
	},
	{
		Name:        "Mayday: leading shift",
		InDate:      mustParse("2021-05-02T12:00:00Z"),
		BankHoliday: false,
	},
	{
		Name:        "Mayday: Midnight shift",
		InDate:      mustParse("2021-05-03T00:00:00Z"),
		BankHoliday: true,
	},
	{
		Name:        "Mayday: Midday shift",
		InDate:      mustParse("2021-05-03T12:00:00Z"),
		BankHoliday: true,
	},
	{
		Name:        "Mayday: following shift",
		InDate:      mustParse("2021-05-04T00:00:00Z"),
		BankHoliday: false,
	},
	{
		Name:        "Mayday: finished",
		InDate:      mustParse("2021-05-04T12:00:00Z"),
		BankHoliday: false,
	},
}

func TestIsBankHoliday(t *testing.T) {
	a := App{
		bankholidays: map[string]struct{}{
			"2021-01-01": {},
			"2021-05-03": {},
		},
	}

	for _, td := range testData {
		if a.isBankholiday(td.InDate) != td.BankHoliday {
			outStr := "NOT be"
			if td.BankHoliday {
				outStr = "be"
			}

			t.Errorf("Failed test: %s", td.Name)
			t.Errorf("Expecting %s to %s a bank holiday.", td.InDate.Format(time.RFC3339), outStr)
		}
	}
}

func mustParse(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		panic(err)
	}

	return t
}
