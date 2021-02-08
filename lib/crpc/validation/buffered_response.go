package validation

import (
	"bytes"
	"net/http"
)

type bufferedResponse struct {
	statusCode int
	body       *bytes.Buffer
}

func newBufferedResponse() *bufferedResponse {
	buf := make([]byte, 0)
	return &bufferedResponse{
		statusCode: http.StatusOK,
		body:       bytes.NewBuffer(buf),
	}
}

func (b bufferedResponse) Header() http.Header {
	// ignore headers
	return nil
}

func (b bufferedResponse) Write(i []byte) (int, error) {
	return b.body.Write(i)
}

func (b bufferedResponse) WriteHeader(statusCode int) {
	b.statusCode = statusCode
}

