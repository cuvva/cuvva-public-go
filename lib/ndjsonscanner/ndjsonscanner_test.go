package ndjsonscanner

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func TestNewBufferScanner(t *testing.T) {
	data, err := os.Open("sample.ndjson")
	assert.NoError(t, err)

	scanner := NewBufferScanner(data)

	defer scanner.Close()

	var actual []TestStruct
	for scanner.Scan() {
		r := TestStruct{}
		err := json.Unmarshal(scanner.Bytes(), &r)
		assert.NoError(t, err)
		actual = append(actual, r)
	}

	assert.NoError(t, scanner.Err())

	expected := []TestStruct{
		{
			A: 1,
			B: "aaa",
		},
		{
			A: 2,
			B: "bbb",
		},
	}

	assert.Equal(t, expected, actual)
}

func TestNewURLScannerWrongHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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

func TestNewScannerEmptyReaderClose(t *testing.T) {
	scanner := NDJSONScanner{
		scanner: nil,
		err:     nil,
		reader:  nil,
	}
	scanner.Close()
}

func TestNewURLScanner(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/nd+json")
		file, err := os.ReadFile("sample.ndjson")
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
		err := json.Unmarshal(scanner.Bytes(), &r)
		assert.NoError(t, err)
		actual = append(actual, r)
	}

	assert.NoError(t, scanner.Err())

	expected := []TestStruct{
		{
			A: 1,
			B: "aaa",
		},
		{
			A: 2,
			B: "bbb",
		},
	}

	assert.Equal(t, expected, actual)
}

func TestNewURLScannerGzip(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/nd+json")
		w.Header().Set("Content-Encoding", "gzip")
		file, err := os.ReadFile("sample.ndjson.gz")
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
		err := json.Unmarshal(scanner.Bytes(), &r)
		assert.NoError(t, err)
		actual = append(actual, r)
	}

	assert.NoError(t, scanner.Err())

	expected := []TestStruct{
		{
			A: 1,
			B: "aaa",
		},
		{
			A: 2,
			B: "bbb",
		},
	}

	assert.Equal(t, expected, actual)
}

func TestNewURLScannerEmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/nd+json")
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
		err := json.Unmarshal(scanner.Bytes(), &r)
		assert.NoError(t, err)
		actual = append(actual, r)
	}

	assert.NoError(t, scanner.Err())

	var expected []TestStruct

	assert.Equal(t, expected, actual)
}

func TestNewBufferScannerScanCallExistsEarly(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/nd+json")
		file, err := os.ReadFile("baddata.ndjson")
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
		err := json.Unmarshal(scanner.Bytes(), &r)
		assert.NoError(t, err)
		callCount++
	}

	assert.NoError(t, scanner.Err())

	expected := 2

	assert.Equal(t, expected, callCount)
}
