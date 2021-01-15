package csvscanner

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	AAA *int    `json:"aaa,string,omitempty"`
	BBB *int    `json:"bbb,string,omitempty"`
	CCC *string `json:"ccc,omitempty"`
}

func TestScannerDoesntSupportBOM(t *testing.T) {
	data, err := os.Open("testdata/bom.csv")
	assert.NoError(t, err)

	_, err = NewCSVScanner(data)
	assert.Error(t, err)
}

func TestNewCSVScanner(t *testing.T) {
	data, err := os.Open("testdata/gooddata.csv")
	assert.NoError(t, err)

	scanner, err := NewCSVScanner(data)
	if err != nil {
		t.Fatal(err)
	}

	defer scanner.Close()

	var actual []TestStruct
	for scanner.Scan() {
		r := TestStruct{}
		_ = json.Unmarshal(scanner.Bytes(), &r)
		actual = append(actual, r)
	}

	assert.NoError(t, scanner.err)

	expected := []TestStruct{
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String("\"\""),
		},
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String("\"\""),
		},
	}

	assert.Equal(t, expected, actual)
}

func TestScannerHandlesWrongHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/error")
		w.Header().Set("Content-Encoding", "gzip")
	}))
	defer ts.Close()

	ctx := context.TODO()

	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	_, err := NewURLScanner(ctx, httpClient, ts.URL)
	assert.Error(t, err)
}

// important to have close method for users to correctly close resources
func TestHasCloseMethod(t *testing.T) {
	scanner := CSVScanner{
		csvReader:  nil,
		err:        nil,
		readCloser: nil,
	}
	scanner.Close()
}

func TestNewURLScanner(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv;charset=utf-8")
		file, err := ioutil.ReadFile("testdata/gooddata.csv")
		if err != nil {
			panic(err)
		}
		_, err = w.Write(file)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	ctx := context.TODO()

	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	scanner, err := NewURLScanner(ctx, httpClient, ts.URL)
	assert.NoError(t, err)

	defer scanner.Close()

	var actual []TestStruct
	for scanner.Scan() {
		r := TestStruct{}
		_ = json.Unmarshal(scanner.Bytes(), &r)
		actual = append(actual, r)
	}

	if err := scanner.Err(); err != nil {
		t.Error(err)
	}

	expected := []TestStruct{
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String("\"\""),
		},
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String("\"\""),
		},
	}

	assert.Equal(t, expected, actual)
}

func TestNewURLScannerGzip(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv;charset=utf-8")
		w.Header().Set("Content-Encoding", "gzip")
		file, err := ioutil.ReadFile("testdata/gooddata.csv.gz")
		if err != nil {
			panic(err)
		}
		_, err = w.Write(file)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	ctx := context.TODO()

	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	scanner, err := NewURLScanner(ctx, httpClient, ts.URL)
	assert.NoError(t, err)

	defer scanner.Close()

	var actual []TestStruct
	for scanner.Scan() {
		r := TestStruct{}
		_ = json.Unmarshal(scanner.Bytes(), &r)
		actual = append(actual, r)
	}

	if err := scanner.Err(); err != nil {
		t.Error(err)
	}

	expected := []TestStruct{
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String("\"\""),
		},
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String("\"\""),
		},
	}

	assert.Equal(t, expected, actual)
}

func TestScannerHandlesEmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv;charset=utf-8")
	}))
	defer ts.Close()

	ctx := context.TODO()

	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	scanner, err := NewURLScanner(ctx, httpClient, ts.URL)
	assert.NoError(t, err)

	defer scanner.Close()

	var actual []TestStruct
	for scanner.Scan() {
		r := TestStruct{}
		_ = json.Unmarshal(scanner.Bytes(), &r)
		actual = append(actual, r)
	}

	assert.NoError(t, scanner.Err())

	var expected []TestStruct

	assert.Equal(t, expected, actual)
}

func TestScanCallExistsEarlyWithBadCSVData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv;charset=utf-8")
		file, err := ioutil.ReadFile("testdata/baddata.csv")
		if err != nil {
			panic(err)
		}
		_, err = w.Write(file)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	ctx := context.TODO()

	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	scanner, err := NewURLScanner(ctx, httpClient, ts.URL)
	assert.NoError(t, err)

	defer scanner.Close()

	callCount := 0
	for scanner.Scan() {
		r := TestStruct{}
		_ = json.Unmarshal(scanner.Bytes(), &r)
		callCount++
	}

	if err := scanner.Err(); err != nil {
		assert.Error(t, err)
	}

	expected := 2

	assert.Equal(t, expected, callCount)
}

func TestScannerErrorsThenNotCSV(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv;charset=utf-8")
		file, err := ioutil.ReadFile("testdata/not_csv.csv")
		if err != nil {
			panic(err)
		}
		_, err = w.Write(file)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	ctx := context.TODO()

	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	scanner, err := NewURLScanner(ctx, httpClient, ts.URL)
	assert.NoError(t, err)

	for scanner.Scan() {
	}

	assert.Error(t, scanner.Err())
}

func TestScannerHandlesEmptyCells(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv;charset=utf-8")
		file, err := ioutil.ReadFile("testdata/emptycells.csv")
		if err != nil {
			panic(err)
		}
		_, err = w.Write(file)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	ctx := context.TODO()

	httpClient := &http.Client{
		Timeout: 50 * time.Second,
	}

	scanner, err := NewURLScanner(ctx, httpClient, ts.URL)
	assert.NoError(t, err)

	defer scanner.Close()

	var actual []TestStruct
	for scanner.Scan() {
		r := TestStruct{}
		_ = json.Unmarshal(scanner.Bytes(), &r)
		actual = append(actual, r)
	}

	if err := scanner.Err(); err != nil {
		t.Error(err)
	}

	expected := []TestStruct{
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String(""),
		},
		{
			AAA: ptr.Int(123),
			BBB: ptr.Int(123),
			CCC: ptr.String(""),
		},
		{
			AAA: ptr.Int(123),
			CCC: ptr.String(""),
		},
	}

	assert.Equal(t, expected, actual)
}
