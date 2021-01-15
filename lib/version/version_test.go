package version

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	tests := []struct {
		Name          string
		App, Revision string
		Expected      string
	}{
		{"Full", "Test", "dev", "Test/dev"},
		{"Truncated", "Test", "e51d8e3c9eca72a41f205a11a2698373bbebb447", "Test/e51d8e3"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			handlerInvoked := false
			w := httptest.NewRecorder()
			r := &http.Request{}
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

			Revision = test.Revision
			Truncated = genTruncated()
			mw := Header(test.App)
			if assert.NotNil(t, mw) {
				hn := mw(next)
				if assert.NotNil(t, hn) {
					hn.ServeHTTP(w, r)

					assert.Equal(t, test.Expected, w.Header().Get("Server"))
					assert.True(t, handlerInvoked, "handler not invoked")
				}
			}
		})
	}
}
