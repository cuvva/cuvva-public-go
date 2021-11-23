package checkmot

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type Client struct {
	*jsonclient.Client
}

// NewClient creates a new vehicle service client.
func NewClient(baseURL, key string) *Client {
	httpClient := &http.Client{
		Transport: &roundTripper{key},
		Timeout:   5 * time.Second,
	}

	return &Client{
		jsonclient.NewClient(baseURL, httpClient),
	}
}

func (c *Client) GetRecordByVRM(ctx context.Context, vrm string) (res *Vehicle, err error) {
	var out []*Vehicle
	if err = c.Do(ctx, "GET", "/trade/vehicles/mot-tests", url.Values{"registration": {vrm}}, nil, &out); err != nil {
		return
	}

	if len(out) == 0 {
		return nil, errors.New(ErrNoResults)
	}

	if len(out) == 1 {
		return out[0], nil
	}

	return attemptToMergeTestsForTheSameVehicle(out)
}

// the API sometimes returns multiple vehicles with the same details even for a specific vehicle search
// merge the list of MOT tests if all vehicle details are the same
func attemptToMergeTestsForTheSameVehicle(out []*Vehicle) (*Vehicle, error) {
	for i := 1; i < len(out); i++ {
		if !isSameVehicle(out[0], out[i]) {
			return nil, errors.New(ErrMultipleVehicles)
		}

		out[0].MOTTests = append(out[0].MOTTests, out[i].MOTTests...)
	}

	return out[0], nil
}

func isSameVehicle(a *Vehicle, b *Vehicle) bool {
	if a.Make != b.Make {
		return false
	}
	if a.Model != b.Model {
		return false
	}
	if a.Registration != b.Registration {
		return false
	}
	if a.FuelType != b.FuelType {
		return false
	}
	if a.PrimaryColour != b.PrimaryColour {
		return false
	}

	return true
}
