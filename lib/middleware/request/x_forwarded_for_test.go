package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCuvvaClientIP(t *testing.T) {
	tests := []struct {
		Name               string
		CuvvaClientIP      string
		ExpectedRemoteAddr string
	}{
		{"NotSet", "", "1.2.3.4"},
		{"Whitespace", "", "1.2.3.4"},
		{"Set", "8.8.4.4", "8.8.4.4"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			handlerInvoked := false
			w := httptest.NewRecorder()
			r := &http.Request{Header: http.Header{CuvvaClientIPHeader: []string{test.CuvvaClientIP}}, RemoteAddr: "1.2.3.4"}
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

			hn := CuvvaClientIP(next)
			if assert.NotNil(t, hn) {
				hn.ServeHTTP(w, r)

				assert.Equal(t, test.ExpectedRemoteAddr, r.RemoteAddr)
				assert.True(t, handlerInvoked, "handler not invoked")
			}
		})
	}
}
