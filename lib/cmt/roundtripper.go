package cmt

import (
	"net/http"
)

type roundTripper struct {
	apiKey string
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Cmt-Api-Key", rt.apiKey)

	return http.DefaultTransport.RoundTrip(req)
}
