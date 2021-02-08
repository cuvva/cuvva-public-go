package validation

import (
	"bytes"
	"net/http"
)

type CopyResponseWriter struct {
	client     http.ResponseWriter
	StatusCode int
	Body       *bytes.Buffer
}

func NewCopyResponseWriter(client http.ResponseWriter) *CopyResponseWriter {
	buf := make([]byte, 0)

	return &CopyResponseWriter{
		client:     client,
		StatusCode: http.StatusOK,
		Body:       bytes.NewBuffer(buf),
	}
}

// We only write headers to the client at the moment because validation doesn't care
func (d *CopyResponseWriter) Header() http.Header {
	return d.client.Header()
}

func (d *CopyResponseWriter) Write(i []byte) (int, error) {
	_, _ = d.Body.Write(i)
	return d.client.Write(i)
}

func (d *CopyResponseWriter) WriteHeader(statusCode int) {
	d.StatusCode = statusCode
	d.client.WriteHeader(statusCode)
}
