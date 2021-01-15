package checkmot

import (
	"encoding/json"
	"time"
)

var loc, _ = time.LoadLocation("Europe/London")

const dateFormatInt = "2006-01-02"
const dateFormatExt = "2006.01.02"
const timeFormat = "2006.01.02 15:04:05"

// Date accepts the DVSA's silly date format (2006.01.02) and converts it to
// ISO8601 while unmarshalling. When marshalling, it remains as ISO8601.
type Date string

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return
	}

	var str string
	if err = json.Unmarshal(data, &str); err != nil {
		return
	}

	if parsed, err := time.ParseInLocation(dateFormatExt, str, loc); err == nil {
		*d = Date(parsed.Format(dateFormatInt))
	}
	return
}

// Time accepts the DVSA's silly time format (2006.01.02 15:04:05) and converts
// it to ISO8601 while unmarshalling. When marshalling, it remains as ISO8601.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return
	}

	var str string
	if err = json.Unmarshal(data, &str); err != nil {
		return
	}

	if parsed, err := time.ParseInLocation(timeFormat, str, loc); err == nil {
		*t = Time{parsed}
	}
	return
}
