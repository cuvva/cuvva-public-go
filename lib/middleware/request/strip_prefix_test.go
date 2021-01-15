package request

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripPrefix(t *testing.T) {
	tests := []struct {
		Name   string
		Path   string
		Prefix string
		Result string
	}{
		{"Empty", "/", "", "/"},
		{"HasPrefix", "/foo/bar", "/foo", "/bar"},
		{"NoPrefix", "/foo/bar", "/baz", "/foo/bar"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var invoked bool

			hn := StripPrefix(test.Prefix)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				invoked = true

				assert.Equal(t, test.Result, r.URL.Path)
			}))

			hn.ServeHTTP(nil, &http.Request{
				URL: &url.URL{
					Path: test.Path,
				},
			})

			assert.True(t, invoked, "handler not invoked")
		})
	}
}
