package sanctionsio

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type Client struct {
	*jsonclient.Client
	apiKey string
}

// NewClient creates a new sanctions.io API client.
func NewClient(baseURL, apiKey string) *Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		jsonclient.NewClient(baseURL, httpClient),
		apiKey,
	}
}

// Search queries the API with options for Name, Sources and Date of Birth
func (c *Client) Search(ctx context.Context, req *SearchRequest) (res *SearchResponse, err error) {
	if req.DateOfBirth != nil && !SearchRequestDOBFormat.MatchString(*req.DateOfBirth) {
		return nil, errors.New("invalid date of birth")
	}

	q := req.URLValues()
	q.Set("api_key", c.apiKey)

	return res, c.Do(ctx, "GET", "/search/", q, nil, &res)
}
