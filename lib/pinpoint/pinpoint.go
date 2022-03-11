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
	ID                    string        `json:"id"`
	Title                 string        `json:"title"`
	EmploymentType        string        `json:"employment_type"`
	Department            JobDepartment `json:"department"`
	ReportingTo           string        `json:"reporting_to"`
	CompensationMinimum   float64       `json:"compensation_minimum"`
	CompensationMaximum   float64       `json:"compensation_maximum"`
	CompensationCurrency  string        `json:"compensation_currency"`
	CompensationFrequency string        `json:"compensation_frequency"`
	CompensationVisible   bool          `json:"compensation_visible"`
	URL                   string        `json:"url"`
	Location              JobLocation   `json:"location"`
}

type JobDepartment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JobLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
