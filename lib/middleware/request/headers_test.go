package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	t.Run("SetHeader", func(t *testing.T) {
		handlerInvoked := false
		w := httptest.NewRecorder()
		r := &http.Request{}
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

		mw := SetHeader("X-Powered-By", "Go!")
		if assert.NotNil(t, mw) {
			hn := mw(next)
			if assert.NotNil(t, hn) {
				hn.ServeHTTP(w, r)

				assert.Equal(t, "Go!", w.Header().Get("X-Powered-By"))
				assert.True(t, handlerInvoked, "handler not invoked")
			}
		}
	})
}
