package csvscanner

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
)

type CSVScanner struct {
	csvReader   *csv.Reader
	ColumnNames []string
	csvRow      []string
	jsonObject  map[string]string
	readCloser  io.ReadCloser
	err         error
}

// setErr records the first error encountered, ignore end of file error
func (cs *CSVScanner) setErr(err error) {
	if cs.err != nil || err == nil || err == io.EOF {
		return
	}

	cs.err = err
}

// err returns the first non-EOF error that was encountered by the CSVScanner
func (cs *CSVScanner) Err() error {
	return cs.err
}

func (cs *CSVScanner) setColumns() {
	cs.Scan()
	cs.ColumnNames = cs.Row()
}

func (cs *CSVScanner) Unmarshal(t interface{}) {
	err := json.Unmarshal(cs.Bytes(), t)
	if err != nil {
		cs.setErr(fmt.Errorf("failed to unmarshal %s: %w", cs.jsonObject, err))
	}
}

func (cs *CSVScanner) Bytes() []byte {
	b, err := json.Marshal(cs.jsonObject)
	if err != nil {
		cs.setErr(err)
		return nil
	}
	return b
}

func (cs *CSVScanner) Close() error {
	if cs.readCloser != nil {
		return cs.readCloser.Close()
	}

	return nil
}

func (cs *CSVScanner) Row() []string {
	return cs.csvRow
}

func (cs *CSVScanner) Scan() bool {
	record, err := cs.csvReader.Read()
	switch err := err; {
	case err == io.EOF:
		return false
	case err != nil:
		cs.setErr(err)
		return false
	}

	cs.csvRow = record

	jsonObject := make(map[string]string)
	for i, columnName := range cs.ColumnNames {
		jsonObject[columnName] = record[i]
	}

	cs.jsonObject = jsonObject

	return true
}

func resolveHTTPClient(httpClient *http.Client) *http.Client {
	if httpClient != nil {
		return httpClient
	}
	return http.DefaultClient
}

func NewURLScanner(ctx context.Context, httpClient *http.Client, url string) (*CSVScanner, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	c := resolveHTTPClient(httpClient)

	// close call on CSVScanner will close the body
	res, err := c.Do(req) //nolint:bodyclose
	if err != nil {
		return nil, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, fmt.Errorf("failed to retrieve url data: %s", res.Status)
	}

	contentType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve content type: %w", err)
	}
	if contentType != "text/csv" {
		return nil, fmt.Errorf("unsupported data format; need 'text/csv', got: %s", contentType)
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

		return NewCSVScanner(gzipReader)
	default:
		return NewCSVScanner(res.Body)
	}
}

func NewCSVScanner(reader io.ReadCloser) (*CSVScanner, error) {
	csvReader := csv.NewReader(reader)
	csvScanner := CSVScanner{
		csvReader:  csvReader,
		readCloser: reader,
	}

	csvScanner.setColumns()

	// using first column header to check for DOM value, allows to avoid buffer ahead scanning complexity
	if len(csvScanner.ColumnNames) > 0 {
		if startsWithBOM(csvScanner.ColumnNames[0]) {
			return nil, errors.New("doesnt support files with BOM")
		}
	}

	return &csvScanner, nil
}
