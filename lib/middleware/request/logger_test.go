package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestResponseWriter(t *testing.T) {
	t.Run("WriteHeader", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: w}

		rw.WriteHeader(http.StatusTeapot)

		assert.Equal(t, http.StatusTeapot, w.Code)
		assert.Equal(t, http.StatusTeapot, rw.Status)
	})

	t.Run("Write", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: w}

		data := []byte("hello")

		n, err := rw.Write(data)
		if assert.NoError(t, err) {
			assert.Equal(t, n, 5)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, http.StatusOK, rw.Status)
			assert.Equal(t, data, w.Body.Bytes())
		}
	})
}

func TestLogger(t *testing.T) {
	tests := []struct {
		Name     string
		Status   int
		Contains string
	}{
		{"Error", http.StatusInternalServerError, "request"},
		{"Warning", http.StatusBadRequest, "request"},
		{"Success", http.StatusOK, "request"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			log := logrus.New().WithField("foo", "bar")

			var buf bytes.Buffer
			log.Logger.Out = &buf

			data := []byte("hello")

			handlerInvoked := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				handlerInvoked = true
				w.WriteHeader(test.Status)
				w.Write(data)
			})

			mw := Logger(log)
			if assert.NotNil(t, mw) {
				fn := mw(next)
				if assert.NotNil(t, fn) {
					w := httptest.NewRecorder()
					r := &http.Request{
						Method:     "GET",
						URL:        &url.URL{Path: "/"},
						Proto:      "HTTP/1.1",
						RemoteAddr: "127.0.0.1",
						Header: http.Header{
							"User-Agent": []string{"FooBar"},
							"Referer":    []string{"FooBar"},
						},
					}

					fn.ServeHTTP(w, r)

					assert.Equal(t, test.Status, w.Code)
					assert.Equal(t, data, w.Body.Bytes())
					assert.True(t, handlerInvoked, "handler not invoked")

					// TODO(jc): compare whole log entry, currently contains timestamp so
					// plain comparison will not work
					assert.Contains(t, buf.String(), test.Contains)
				}
			}
		})
	}
}
