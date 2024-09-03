package appsflyer

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/ptr"
	"github.com/stretchr/testify/assert"
)

func TestFilterNextKey(t *testing.T) {
	cases := []struct {
		Test             string
		S3Keys           []string
		PresentS3Key     *string
		ExpectedNewS3Key *string
		ExpectedError    bool
	}{
		{
			Test: "will return first value from list if cursor is empty",
			S3Keys: []string{
				"a/1/2/h=4/some_file.gz",
				"a/1/2/h=3/some_file.gz",
			},
			ExpectedNewS3Key: ptr.String("a/1/2/h=3/some_file.gz"),
			ExpectedError:    false,
		},
		{
			Test: "will correctly order and pick up first value from s3 files list",
			S3Keys: []string{
				"a/1/2/h=10/some_file.gz",
				"a/1/2/h=3/some_file.gz",
			},
			ExpectedNewS3Key: ptr.String("a/1/2/h=3/some_file.gz"),
			ExpectedError:    false,
		},
		{
			Test:          "has to throw error if s3 keys list is empty",
			ExpectedError: true,
		},
		{
			Test: "correctly identifies that h=3 is less than h=10",
			S3Keys: []string{
				"a/1/2/h=10/some_file.gz",
				"a/1/2/h=3/some_file.gz",
				"a/1/2/h=11/some_file.gz",
			},
			PresentS3Key:     ptr.String("a/1/2/h=10/some_file.gz"),
			ExpectedNewS3Key: ptr.String("a/1/2/h=11/some_file.gz"),
			ExpectedError:    false,
		},
		{
			Test: "returns nil and error if there is no new files",
			S3Keys: []string{
				"85c7-acc-VNBuiQLC-85c7/data-locker-hourly/t=sessions/dt=2020-06-19/h=17/part-00000.gz",
				"85c7-acc-VNBuiQLC-85c7/data-locker-hourly/t=sessions/dt=2020-06-19/h=18/part-00000.gz",
				"85c7-acc-VNBuiQLC-85c7/data-locker-hourly/t=sessions/dt=2020-06-19/h=19/part-00000.gz",
				"85c7-acc-VNBuiQLC-85c7/data-locker-hourly/t=sessions/dt=2020-06-19/h=20/part-00000.gz",
				"85c7-acc-VNBuiQLC-85c7/data-locker-hourly/t=sessions/dt=2020-06-19/h=21/part-00000.gz",
			},
			PresentS3Key:     ptr.String("85c7-acc-VNBuiQLC-85c7/data-locker-hourly/t=sessions/dt=2020-06-19/h=21/part-00000.gz"),
			ExpectedError:    true,
			ExpectedNewS3Key: nil,
		},
	}
	for _, c := range cases {
		s3Key, err := filterNextS3Key(c.S3Keys, c.PresentS3Key)
		t.Log(c.Test)
		if c.ExpectedError {
			assert.Error(t, err)
			assert.Nil(t, c.ExpectedNewS3Key)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, s3Key, c.ExpectedNewS3Key)
		}
	}
}
