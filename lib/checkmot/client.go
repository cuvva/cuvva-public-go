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

	if len(out) != 1 {
		return nil, errors.New(ErrUnexpectedResultCount)
	}

	return out[0], nil
}
