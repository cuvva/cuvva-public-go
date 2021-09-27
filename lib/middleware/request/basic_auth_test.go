package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	server := NewBasicAuth("username", "supersecretpassword")

	t.Run("authentication passes", func(t *testing.T) {
		handlerInvoked := false

		w := httptest.NewRecorder()
		r := &http.Request{Header: http.Header{"Authorization": []string{"Basic dXNlcm5hbWU6c3VwZXJzZWNyZXRwYXNzd29yZA=="}}}
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

		hn := NewBasicAuthMiddleware(server.CheckAuth)(next)
		if assert.NotNil(t, hn) {
			hn.ServeHTTP(w, r)

			assert.Equal(t, 200, w.Code)
			assert.True(t, handlerInvoked, "handler not invoked")
		}
	})

	t.Run("authentication rejected", func(t *testing.T) {
		handlerInvoked := false

		w := httptest.NewRecorder()
		r := &http.Request{Header: http.Header{"Authorization": []string{"Basic dXNlcm5hbWU6aGFja2VybWFu"}}}
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

		hn := NewBasicAuthMiddleware(server.CheckAuth)(next)
		if assert.NotNil(t, hn) {
			hn.ServeHTTP(w, r)

			assert.Equal(t, 401, w.Code)
			assert.False(t, handlerInvoked, "handler not invoked")
		}
	})
}
