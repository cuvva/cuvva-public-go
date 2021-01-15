package dvlaves

import (
	"context"
	"net/http"
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

func (c *Client) GetVehicleByVRM(ctx context.Context, vrm string) (res *Vehicle, err error) {
	return res, c.Do(ctx, "POST", "vehicle-enquiry/v1/vehicles", nil, &VESVRMRequest{RegistrationNumber: vrm}, &res)
}
