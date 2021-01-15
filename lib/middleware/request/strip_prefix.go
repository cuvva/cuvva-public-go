package request

import (
	"net/http"
	"strings"
)

// StripPrefix is HTTP Middleware that will strip a prefix
// from a URL before hitting an endpoint.
func StripPrefix(prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, prefix)

			next.ServeHTTP(w, r)
		})
	}
}
