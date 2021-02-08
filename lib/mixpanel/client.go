package mixpanel

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
)

type Client struct {
	*http.Client
}

type BasicAuthTransport struct {
	APISecret string
}

func (ba BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(ba.APISecret, "")

	return http.DefaultTransport.RoundTrip(req)
}

func NewClient(apiSecret string) *Client {
	basicAuthTransport := BasicAuthTransport{APISecret: apiSecret}
	httpClient := Client{&http.Client{
		Transport: basicAuthTransport,
	}}

	return &httpClient
}

func (c Client) Export(ctx context.Context, fromDate, toDate time.Time) (ers *ExportResultScanner, err error) {
	// https://developer.mixpanel.com/docs/exporting-raw-data
	urlTemplate := `https://data.mixpanel.com/api/2.0/export/?from_date=%s&to_date=%s`
	url := fmt.Sprintf(urlTemplate,
		fromDate.Format("2006-01-02"),
		toDate.Format("2006-01-02"),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, cher.New(cher.TooManyRequests, nil)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, cher.New(cher.Unknown, cher.M{"error": err})
		}

		return nil, cher.New(cher.Unknown, cher.M{"response_status": resp.Status,
			"response_body": string(body)})
	}

	bns := bufio.NewScanner(resp.Body)

	return NewExportResultScanner(*bns), nil
}
