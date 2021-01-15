package ukvd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

const apiNullItems = "1"
const apiVersion = "2"

type Client struct {
	client *jsonclient.Client
	key    string
}

func NewClient(baseURL, key string) *Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		jsonclient.NewClient(baseURL, httpClient),
		key,
	}
}

func (c *Client) GetFuelPriceData(ctx context.Context, latitude float64, longitude float64) (res *FuelPriceData, err error) {

	params := url.Values{
		"v":             []string{apiVersion},
		"api_nullitems": []string{apiNullItems},
		"auth_apikey":   []string{c.key},
		"key_latitude":  []string{fmt.Sprintf("%f", latitude)},
		"key_longitude": []string{fmt.Sprintf("%f", longitude)},
	}

	return res, c.client.Do(ctx, "GET", "/api/datapackage/FuelPriceData", params, nil, &res)
}
