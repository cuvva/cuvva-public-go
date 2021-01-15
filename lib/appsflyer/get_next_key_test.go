package appsflyer

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/stretchr/testify/assert"
)

func TestScanMakeS3KeyOrderable(t *testing.T) {
	cases := []struct {
		Test     string
		S3Key    string
		Expected *string
	}{
		{
			Test:     "h=25 to h=25",
			S3Key:    "1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/h=25/part-00000.gz",
			Expected: ptr.String("1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/h=25/part-00000.gz"),
		},
		{
			Test:     "h=2 to h=02",
			S3Key:    "1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/h=2/part-00000.gz",
			Expected: ptr.String("1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/h=02/part-00000.gz"),
		},
		{
			Test:     "h=20 to h=20",
			S3Key:    "1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/h=20/part-00000.gz",
			Expected: ptr.String("1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/h=20/part-00000.gz"),
		},
	}

	for _, c := range cases {
		actual := makeS3KeyOrderable(c.S3Key)
		assert.Equal(t, *c.Expected, actual)
	}
}

func TestGetApsFlyerReportDatePart(t *testing.T) {
	cases := []struct {
		Test          string
		S3Key         string
		Expected      string
		ExpectedError bool
	}{
		{
			Test:          "finds /h=",
			S3Key:         "1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/h=20/part-00000.gz",
			Expected:      "1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17",
			ExpectedError: false,
		},
		{
			Test:          "missing /h=",
			S3Key:         "1cd7-acc-H8JAFbQb-1cd7/data-locker-hourly/t=clicks/dt=2020-04-17/aaaa=20/part-00000.gz",
			Expected:      "",
			ExpectedError: true,
		},
		{
			Test:          "fails on nil input",
			Expected:      "",
			ExpectedError: true,
		},
	}

	for _, c := range cases {
		actual, err := trimReportPathToDate(&c.S3Key)
		if c.ExpectedError {
			assert.Error(t, err)
			assert.Equal(t, c.Expected, actual)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.Expected, actual)
		}
	}
}
