package cmt

import (
	"context"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type Client interface {
	RequestAuthorizationCode(ctx context.Context, req *AuthorizationCodeRequest) (*AuthorizationCodeResponse, error)
}

type HTTPClient struct {
	*jsonclient.Client
}

func NewClient(baseURL string, apiKey string) *HTTPClient {
	httpClient := &http.Client{
		Transport: &roundTripper{apiKey: apiKey},
		Timeout:   5 * time.Second,
	}

	return &HTTPClient{
		jsonclient.NewClient(baseURL, httpClient),
	}
}

func (c *HTTPClient) RequestAuthorizationCode(ctx context.Context, req *AuthorizationCodeRequest) (res *AuthorizationCodeResponse, err error) {
	return res, c.Do(ctx, "POST", "/v4/request_auth_code", nil, req, &res)
}
