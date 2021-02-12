package app

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/cuvva/cuvva-public-go/tools/oncallcalc/lib/govuk"
)

// App contains all the business logic for this tool
type App struct {
	pagerduty    *pagerduty.Client
	govUK        *govuk.Client
	bankholidays map[string]struct{}
}

// New helps initialise a new App
func New(client *pagerduty.Client, govUK *govuk.Client) *App {
	return &App{client, govUK, make(map[string]struct{})}
}
