package dvlaves

import (
	"net/http"
)

type roundTripper struct {
	apiKey string
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-api-key", rt.apiKey)

	return http.DefaultTransport.RoundTrip(req)
}
