package restbase

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestREST(t *testing.T) {
	t.Run("Send", func(t *testing.T) {
		tests := []struct {
			Name    string
			Request *http.Request
			Body    interface{}

			Headers    http.Header
			Bytes      []byte
			StatusCode int
		}{
			{"NoContent", &http.Request{}, nil, http.Header{}, nil, http.StatusNoContent},
			{
				"JSON", &http.Request{Header: http.Header{"Accept": []string{"application/json"}}},
				"foo", http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
				[]byte(`"foo"` + "\n"), http.StatusOK,
			},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				w := httptest.NewRecorder()

				Send(test.Request.Context(), w, test.Body)

				assert.Equal(t, test.StatusCode, w.Code)
				assert.Equal(t, test.Bytes, w.Body.Bytes())
				assert.Equal(t, test.Headers, w.Header())
			})
		}
	})
}
