package ndjsonscanner

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
)

const MaxScanTokenSize = 2000 * 1000

type NDJSONScanner struct {
	scanner       *bufio.Scanner
	err           error
	reader        io.ReadCloser
	contentLength int64
}

func (ns *NDJSONScanner) setErr(err error) {
	ns.err = err
}

func (ns *NDJSONScanner) Err() error {
	return ns.scanner.Err()
}

func (ns *NDJSONScanner) Bytes() []byte {
	return ns.scanner.Bytes()
}

func (ns *NDJSONScanner) Close() error {
	if ns.reader != nil {
		return ns.reader.Close()
	}

	return nil
}

func (ns *NDJSONScanner) ContentLength() int64 {
	return ns.contentLength
}

func (ns *NDJSONScanner) Scan() bool {
	if ns.scanner.Err() != nil {
		return false
	}

	proceed := ns.scanner.Scan()

	var a map[string]interface{}
	err := json.Unmarshal(ns.scanner.Bytes(), &a)
	if err != nil {
		ns.setErr(err)
		return false
	}

	return proceed
}

func resolveHTTPClient(httpClient *http.Client) *http.Client {
	if httpClient != nil {
		return httpClient
	}
	return http.DefaultClient
}

func NewURLScanner(ctx context.Context, httpClient *http.Client, url string) (*NDJSONScanner, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	c := resolveHTTPClient(httpClient)

	// Close call on NDJSONScanner will close the body.
	res, err := c.Do(req) //nolint:bodyclose
	if err != nil {
		return nil, err
	}

	if res.Body == nil {
		return NewBufferScanner(nil), nil
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, fmt.Errorf("failed to retrieve url data: %s", res.Status)
	}

	contentType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve content type: %w", err)
	}
	if contentType != "application/nd+json" {
		return nil, fmt.Errorf("unsupported data format; need 'application/nd+json', got: %s", contentType)
	}

	contentEncoding, _, err := mime.ParseMediaType(res.Header.Get("Content-Encoding"))
	if err != nil && contentEncoding != "" {
		return nil, fmt.Errorf("failed to retrieve content encoding: %w", err)
	}

	switch contentEncoding {
	case "gzip":
		gzipReader, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}

		nbs := NewBufferScanner(gzipReader)
		nbs.contentLength = res.ContentLength

		return nbs, nil
	default:
		nbs := NewBufferScanner(res.Body)
		nbs.contentLength = res.ContentLength

		return nbs, nil
	}
}

func NewBufferScanner(reader io.ReadCloser) *NDJSONScanner {
	if reader != nil {
		scanner := bufio.NewScanner(reader)
		buf := make([]byte, MaxScanTokenSize)
		scanner.Buffer(buf, MaxScanTokenSize)

		return &NDJSONScanner{
			scanner: scanner,
			reader:  reader,
		}
	}

	return &NDJSONScanner{}
}
