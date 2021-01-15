package appsflyer

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cuvva/cuvva-public-go/lib/csvscanner"
)

type Scanner struct {
	csvReader   *csvscanner.CSVScanner
	objectType  interface{}
	typedRecord interface{}
	err         error
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) Bytes() []byte {
	bytes, err := json.Marshal(s.typedRecord)
	if err != nil {
		s.err = err
	}

	return bytes
}

// Scan advances the CSVScanner to the next row, which will then be
// converted to go type and available through the Bytes method. It returns false when the
// scan stops, either by reaching the end of the input or an error.
// After Scan returns false, the Err method will return any error that
// occurred during scanning, except that if it was io.EOF, Err
// will return nil.
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}
	ok := s.csvReader.Scan()
	if !ok {
		s.err = s.csvReader.Err()
		return ok
	}

	rawRecord := make(map[string]interface{})

	for i, columnName := range s.csvReader.ColumnNames {
		rawRecord[columnName] = s.csvReader.Row()[i]
	}

	marshalledRawRecord, err := json.Marshal(rawRecord)
	if err != nil {
		s.err = err
		return false
	}

	typedRecord := s.objectType

	err = json.Unmarshal(marshalledRawRecord, &typedRecord)
	if err != nil {
		s.err = err
		return false
	}

	s.typedRecord = typedRecord

	return true
}

type ParseCSVRowFunc func(data []string, columnNames []string, objectType interface{}) (interface{}, error)

func NewScanner(s3ObjectBody io.Reader, objectType interface{}) (*Scanner, error) {
	gzipData, err := gzip.NewReader(s3ObjectBody)
	if err != nil {
		return nil, fmt.Errorf("failed to creat gzip reader: %w", err)
	}

	csvScanner, err := csvscanner.NewCSVScanner(gzipData)
	if err != nil {
		return nil, fmt.Errorf("failed to create csv scanner: %w", err)
	}

	return &Scanner{
		csvReader:  csvScanner,
		objectType: objectType,
	}, nil
}
