package appsflyer

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
)

var NoData = errors.New("no data")
var NoNewData = errors.New("no new data")

func makeS3KeyOrderable(s3Key string) string {
	// Suffixing single hour digits with 0.
	for i := 1; i < 10; i++ {
		original := fmt.Sprintf("/h=%d/", i)
		replacement := fmt.Sprintf("/h=0%d/", i)
		s3Key = strings.Replace(s3Key, original, replacement, 1)
	}

	s3Key = strings.Replace(s3Key, "/h=late/", "/h=25/", 1)

	return s3Key
}

func trimReportPathToDate(s3Key *string) (string, error) {
	if s3Key == nil {
		return "", nil
	}

	if strings.Contains(*s3Key, "/h=") {
		return strings.Split(*s3Key, "/h=")[0], nil
	}

	return "", fmt.Errorf("unsupported report path: %s", *s3Key)
}

func filterNextS3Key(s3Keys []string, presentS3Key *string) (*string, error) {
	if len(s3Keys) == 0 {
		return nil, errors.New("s3keys list is empty")
	}

	// h=3, and h=10 strings are not suitable for correct ordering
	keyMapping := make(map[string]string)
	orderableKeys := make([]string, 0, len(s3Keys))

	for _, k := range s3Keys {
		orderableKey := makeS3KeyOrderable(k)
		orderableKeys = append(orderableKeys, orderableKey)
		keyMapping[orderableKey] = k
	}

	sort.Strings(orderableKeys)

	if presentS3Key == nil {
		o := keyMapping[orderableKeys[0]]
		return &o, nil
	}

	orderablePresentS3Key := makeS3KeyOrderable(*presentS3Key)

	for _, oKey := range orderableKeys {
		if oKey > orderablePresentS3Key {
			o := keyMapping[oKey]
			return &o, nil
		}
	}

	return nil, NoNewData
}

func (c Client) getNextKey(ctx context.Context, s3Key *string, reportPath string) (*string, error) {
	startAfter, err := trimReportPathToDate(s3Key)
	if err != nil {
		return nil, err
	}

	// Getting first batch of keys
	objects, err := c.s3Client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Prefix:     &reportPath,
		StartAfter: &startAfter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// If list is empty try again, there will be at least one SUCCESS file
	// https://support.appsflyer.com/hc/en-us/articles/360000877538-Data-Locker-V2-0-high-volume-multi-app-raw-data-reporting
	if len(objects.Contents) == 0 {
		return nil, NoData
	}

	var s3Keys []string
	for _, o := range objects.Contents {
		// Code is only extracting .gz files
		if strings.HasSuffix(*o.Key, ".gz") {
			s3Keys = append(s3Keys, *o.Key)
		}
	}

	return filterNextS3Key(s3Keys, s3Key)
}
