package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

// PDShift is a PagerDuty representation of an on-call shift
type PDShift struct {
	Email string
	Start time.Time
	End   time.Time
}

// ListPagerDutyShifts returns a slice of shifts pagerduty is tracking for a schedule
func (a *App) ListPagerDutyShifts(scheduleIDs []string, start, end time.Time) ([]*PDShift, interface{}, error) {
	var shifts []*PDShift

	fmt.Println("Requesting on call dates for ", start.Format(time.RFC3339), " until ", end.Format(time.RFC3339))

	// @TODO context should be injected from the cobra entry point
	res, err := a.pagerduty.ListOnCallsWithContext(context.Background(), pagerduty.ListOnCallOptions{
		Includes:    []string{"users"},
		ScheduleIDs: scheduleIDs,
		Since:       start.Format(time.RFC3339),
		Until:       end.Format(time.RFC3339),
		Limit:       100, // TODO implement pagination if we require more values
	})

	if res.More {
		return nil, nil, errors.New("Pagination is required")
	}

	if err != nil {
		return nil, res, err
	}

	for _, pgShift := range res.OnCalls {
		if pgShift.Start == "" || pgShift.End == "" {
			return nil, pgShift, errors.New("malformed shift start or end time")
		}

		start := MustParse(time.RFC3339, pgShift.Start)
		end := MustParse(time.RFC3339, pgShift.End)

		// if you weren't on the shift for at least 1 hour, we're skipping you
		if end.Sub(start).Seconds() < (60 * time.Minute).Seconds() {
			fmt.Println("Skipping short shift for [ ", pgShift.User.ID, " ] as duration was ", end.Sub(start).Minutes(), " minutes on ", pgShift.Start)
			continue
		}

		shifts = append(shifts, &PDShift{
			Email: pgShift.User.Email,
			Start: start,
			End:   end,
		})
	}

	return shifts, nil, nil
}
