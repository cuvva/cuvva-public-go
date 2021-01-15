package givefood

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type Client struct {
	client *jsonclient.Client
}

func NewClient(baseURL string) *Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		jsonclient.NewClient(baseURL, httpClient),
	}
}

func (c *Client) GetFoodBanks(ctx context.Context) (res []*FoodBank, err error) {
	return res, c.client.Do(ctx, "GET", "/api/1/foodbanks/", nil, nil, &res)
}

func (c *Client) GetFoodBankBySlug(ctx context.Context, slug string) (res *FoodBank, err error) {
	return res, c.client.Do(ctx, "GET", fmt.Sprintf("/api/1/foodbank/%s", slug), nil, nil, &res)
}
