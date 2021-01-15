package restbase

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		Name string
		Src  interface{}

		StatusCode int
		Out        []byte
		MIMEType   string
		Error      error
	}{
		{"DefaultEncodingJSON", "foo", http.StatusOK, []byte(`"foo"` + "\n"), "application/json; charset=utf-8", nil},
		{"NoAccept", "foo", http.StatusOK, []byte(`"foo"` + "\n"), "application/json; charset=utf-8", nil},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := encode(w, test.Src)

			if test.Error == nil {
				if assert.NoError(t, err) {
					assert.Equal(t, w.Body.Bytes(), test.Out)
					assert.Equal(t, w.Code, test.StatusCode)
					assert.Equal(t, w.Header(), http.Header{"Content-Type": []string{test.MIMEType}})
				}
			} else {
				assert.Equal(t, test.Error, err)
			}
		})
	}
}
