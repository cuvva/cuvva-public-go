package pinpoint

import (
	"context"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type JobsResponse struct {
	Data []Job `json:"data"`
}

type Job struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes JobAttributes `json:"attributes"`
	Links      JobLinks      `json:"links"`
}

type JobAttributes struct {
	DepartmentID   int    `json:"department_id"`
	ID             int    `json:"id"`
	Title          string `json:"title"`
	LocationID     int    `json:"location_id"`
	EmploymentType string `json:"employment_type"`
	LocationName   string `json:"location_name"`
	DepartmentName string `json:"department_name"`
}

type JobLinks struct {
	ShowPath string `json:"show_path"`
}

type Client struct {
	*jsonclient.Client
}

func New(baseURL string) *Client {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{jsonclient.NewClient(baseURL, httpClient)}
}

func (c *Client) GetJobs(ctx context.Context) ([]Job, error) {
	var response *JobsResponse

	if err := c.Do(ctx, "GET", "jobs.json", nil, nil, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
