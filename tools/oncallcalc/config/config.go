package config

import (
	"fmt"
	"os"

	"github.com/PagerDuty/go-pagerduty"
)

// PagerDutyAPIEnv is the name of the environment variable that holds the PagerDuty API authToken
const PagerDutyAPIEnv = "PAGERDUTY_API"

// BuildPDClient builds a new pagerduty client from env variables
func BuildPDClient() (*pagerduty.Client, error) {
	authToken, ok := os.LookupEnv(PagerDutyAPIEnv)
	if !ok {
		return nil, fmt.Errorf("missing env var: %s", PagerDutyAPIEnv)
	}

	client := pagerduty.NewClient(authToken)

	return client, nil
}
