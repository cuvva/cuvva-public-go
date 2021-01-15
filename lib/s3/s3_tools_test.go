package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceUnsafeCharacters(t *testing.T) {
	cases := []struct {
		Test                 string
		S3Key                string
		ReplacementCharacter string
		ExpectedNewS3Key     string
		ExpectedError        bool
	}{
		{
			Test:                 "replaces bad characters",
			S3Key:                ":A:s:09!-:_.*'():::&",
			ReplacementCharacter: "-",
			ExpectedNewS3Key:     "-A-s-09!--_.*'()----",
			ExpectedError:        false,
		},
		{
			Test:                 "leaves good key unchanged",
			S3Key:                "a-b-c-d-0-9",
			ReplacementCharacter: "-",
			ExpectedNewS3Key:     "a-b-c-d-0-9",
			ExpectedError:        false,
		},
		{
			Test:                 "fails if replacement key is not safe",
			S3Key:                "a",
			ReplacementCharacter: "%",
			ExpectedError:        true,
		},
		{
			Test:                 "replace bad characters of actual s3 key",
			S3Key:                "warehouse/quote/data_extraction_service/ExtractQuoteData/1/2020-02-18/2016-09-06T11:37:39.700Z|50f21340-7426-11e6-b35b-c357a05a7a26.ndjson.csv",
			ExpectedNewS3Key:     "warehouse/quote/data_extraction_service/ExtractQuoteData/1/2020-02-18/2016-09-06T11-37-39.700Z-50f21340-7426-11e6-b35b-c357a05a7a26.ndjson.csv",
			ReplacementCharacter: "-",
			ExpectedError:        false,
		},
	}
	for _, c := range cases {
		actual, err := ReplaceUnsafeKeyCharacters(c.S3Key, c.ReplacementCharacter)
		if c.ExpectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.ExpectedNewS3Key, *actual)
		}
	}
}
