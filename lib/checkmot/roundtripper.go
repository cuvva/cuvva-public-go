package checkmot

import (
	"net/http"
)

type roundTripper struct {
	apiKey    string
	transport http.RoundTripper
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json+v4")
	req.Header.Set("X-Api-Key", rt.apiKey)

	return rt.transport.RoundTrip(req)
}
