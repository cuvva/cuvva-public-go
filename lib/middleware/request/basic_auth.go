package request

import (
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

type BasicAuth struct {
	username string
	password string
}

func NewBasicAuthMiddleware(check func(user, pass string) bool) func(fn http.Handler) http.Handler {
	return func(fn http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, _ := r.BasicAuth()
			if !check(user, pass) {
				JSONError(w, cher.New(cher.Unauthorized, nil))
				return
			}
			fn.ServeHTTP(w, r)
		})
	}
}

func NewBasicAuth(username string, password string) *BasicAuth {
	if username == "" || password == "" {
		panic("username and password required for basic auth")
	}

	return &BasicAuth{username: username, password: password}
}

func (w *BasicAuth) CheckAuth(user, pass string) bool {
	return user == w.username && pass == w.password
}
