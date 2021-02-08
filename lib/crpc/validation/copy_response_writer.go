package validation

import "net/http"

type copyResponseWriter struct {
	client       http.ResponseWriter
	responseCopy *bufferedResponse
}

func newCopyResponseWriter(client http.ResponseWriter) *copyResponseWriter {
	return &copyResponseWriter{
		client:       client,
		responseCopy: newBufferedResponse(),
	}
}

// We only write headers to the client at the moment because validation doesn't care
func (d *copyResponseWriter) Header() http.Header {
	return d.client.Header()
}

func (d *copyResponseWriter) Write(i []byte) (int, error) {
	_, _ = d.responseCopy.Write(i)
	return d.client.Write(i)
}

func (d *copyResponseWriter) WriteHeader(statusCode int) {
	d.responseCopy.WriteHeader(statusCode)
	d.client.WriteHeader(statusCode)
}

