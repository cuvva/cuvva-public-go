package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	t.Run("Context", func(t *testing.T) {
		const testRequestID = "test"

		t.Run("GetRequestID", func(t *testing.T) {
			t.Run("Empty", func(t *testing.T) {
				r := &http.Request{}

				requestID := GetRequestID(r)
				assert.Empty(t, requestID)
			})

			t.Run("Context", func(t *testing.T) {

				r := &http.Request{}
				r = r.WithContext(context.WithValue(r.Context(), RequestIDKey, testRequestID))

				requestID := GetRequestID(r)
				assert.Equal(t, testRequestID, requestID)
			})
		})

		t.Run("SetRequestID", func(t *testing.T) {
			r := &http.Request{}

			r = SetRequestID(r, testRequestID)

			requestID := r.Context().Value(RequestIDKey).(string)
			assert.Equal(t, testRequestID, requestID)
		})
	})

	t.Run("Middleware", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			handlerInvoked := false

			w := httptest.NewRecorder()
			r := &http.Request{Header: make(http.Header)}
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

			hn := RequestID(next)
			if assert.NotNil(t, hn) {
				hn.ServeHTTP(w, r)

				assert.NotEmpty(t, w.Header().Get("Request-ID"))
				assert.True(t, handlerInvoked, "handler not invoked")
			}
		})

		t.Run("Existing", func(t *testing.T) {
			handlerInvoked := false

			w := httptest.NewRecorder()
			r := &http.Request{Header: http.Header{"Request-Id": []string{"test"}}}
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

			hn := RequestID(next)
			if assert.NotNil(t, hn) {
				hn.ServeHTTP(w, r)

				assert.Equal(t, "test", w.Header().Get("Request-Id"))
				assert.True(t, handlerInvoked, "handler not invoked")
			}
		})
	})
}
