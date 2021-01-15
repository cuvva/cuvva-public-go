package app

import (
	"github.com/PagerDuty/go-pagerduty"
)

// App contains all the business logic for this tool
type App struct {
	pagerduty *pagerduty.Client
}

// New helps initialise a new App
func New(client *pagerduty.Client) *App {
	return &App{client}
}
