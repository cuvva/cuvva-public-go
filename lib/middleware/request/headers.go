package request

import (
	"net/http"
)

// SetHeader returns a middleware handler that injects the given key/value as a HTTP header.
//
// NOTE: It will overwrite any existing header with that name. Subsequent middleware/handlers
// may overwrite/append to the value set by this middleware.
func SetHeader(key, value string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)

			next.ServeHTTP(w, r)
		})
	}
}
