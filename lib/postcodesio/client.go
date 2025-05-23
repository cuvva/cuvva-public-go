package postcodesio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
	"github.com/cuvva/cuvva-public-go/lib/version"
)

// Postcode is struct to contain the postcode information given in the API response
type Postcode struct {
	Postcode  string  `json:"postcode"`
	Area      string  `json:"area"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// UnmarshalJSON is a custom unmarshaler for the Postcode struct
func (p *Postcode) UnmarshalJSON(data []byte) error {
	type Alias Postcode
	aux := &struct {
		Area string `json:"parish"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal postcode: %w", err)
	}

	area := aux.Area
	if strings.Contains(area, ", unparished area") {
		area = strings.Split(area, ", unparished area")[0]
	}

	p.Area = area

	return nil
}

// DefaultUserAgent is the default user agent to use for the lib if no other
// user agent is given
var DefaultUserAgent = "postcodesio/" + version.Truncated + " (+https://cuvva.com)"

// DefaultBaseURL is the default host for postcodes.io
const DefaultBaseURL = "https://api.postcodes.io"

// Service interface contains all available, exposed methods of postcodes.io
type Service interface {
	Geocode(ctx context.Context, postcode string) (*Postcode, error)
	ReverseGeocode(ctx context.Context, latitude, longitude float64) (*Postcode, error)
}

// Client is the base struct for the methods to be attached to
type Client struct {
	*jsonclient.Client
}

// FailoverClient contains many clients and will attempt to execute and
// client operations on them in order until the first non-error response
// is encountered.
type FailoverClient struct {
	clients []*Client
}

// New generates the client struct with populated net/http client
func New(baseURL string) *Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	jcc := jsonclient.NewClient(baseURL, httpClient)

	jcc.UserAgent = DefaultUserAgent

	return &Client{jcc}
}

// NewFailoverClient returns a FailoverClient instance with the provided clients
func NewFailoverClient(clients ...*Client) (*FailoverClient, error) {
	if len(clients) == 0 {
		return nil, errors.New("not enough clients")
	}

	return &FailoverClient{clients}, nil
}

// ReverseGeocode returns a set (or no) postcodes that exist within a long/lat
func (c *Client) ReverseGeocode(ctx context.Context, latitude, longitude float64) (*Postcode, error) {
	params := url.Values{
		"lat":        []string{strconv.FormatFloat(latitude, 'f', -1, 64)},
		"lon":        []string{strconv.FormatFloat(longitude, 'f', -1, 64)},
		"limit":      []string{"1"},
		"radius":     []string{"20000"},
		"wideSearch": []string{"true"},
	}

	var res struct {
		Status int         `json:"status"`
		Result []*Postcode `json:"result"`
	}

	err := c.Do(ctx, "GET", "/postcodes", params, nil, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Result) == 0 {
		return nil, nil
	}

	return res.Result[0], nil
}

func (c *Client) Geocode(ctx context.Context, postcode string) (*Postcode, error) {
	var res struct {
		Status int       `json:"status"`
		Result *Postcode `json:"result"`
	}

	err := c.Do(ctx, "GET", "/postcodes/"+postcode, nil, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

// ReverseGeocode returns a set (or no) postcodes that exist within a long/lat
func (fc *FailoverClient) ReverseGeocode(ctx context.Context, latitude, longitude float64) (*Postcode, error) {
	var errors []cher.E

	for _, cli := range fc.clients {
		pc, err := cli.ReverseGeocode(ctx, latitude, longitude)
		if err == nil {
			return pc, nil
		}

		cErr := cher.New("postcodes_request_failed", cher.M{
			"message": err.Error(),
		})

		errors = append(errors, cErr)
	}

	return nil, cher.New("postcode_error", cher.M{
		"latitude":  latitude,
		"longitude": longitude,
	}, errors...)
}

func (fc *FailoverClient) Geocode(ctx context.Context, postcode string) (*Postcode, error) {
	var errors []cher.E

	for _, cli := range fc.clients {
		pc, err := cli.Geocode(ctx, postcode)
		if err == nil {
			return pc, nil
		}

		cErr := cher.New("postcodes_request_failed", cher.M{
			"message": err.Error(),
		})

		errors = append(errors, cErr)
	}

	return nil, cher.New("postcode_error", cher.M{
		"postcode": postcode,
	}, errors...)
}

// interface guards
var (
	_ Service = (*Client)(nil)
	_ Service = (*FailoverClient)(nil)
)
