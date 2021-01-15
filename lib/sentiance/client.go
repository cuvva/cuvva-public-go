package sentiance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
)

type Client interface {
	UserLink(ctx context.Context, req *UserLinkRequest) (*UserLinkResponse, error)
	GetEventAtStartDate(ctx context.Context, req *GetEventAtStartDateRequest) (*GetEventAtStartDateResponse, error)
}

type HTTPClient struct {
	*jsonclient.Client
}

func NewClient(token string) *HTTPClient {
	httpClient := &http.Client{
		Transport: jsonclient.NewAuthenticatedRoundTripper(nil, "Bearer", token),
		Timeout:   5 * time.Second,
	}

	return &HTTPClient{
		jsonclient.NewClient("https://api.sentiance.com/", httpClient),
	}
}

func (c *HTTPClient) UserLink(ctx context.Context, req *UserLinkRequest) (res *UserLinkResponse, err error) {
	return res, c.Do(ctx, "POST", fmt.Sprintf("/v2/users/%s/link", req.InstallID), nil, UserLink{ExternalID: req.ExternalID}, &res)
}

func (c *HTTPClient) GetEventAtStartDate(ctx context.Context, req *GetEventAtStartDateRequest) (res *GetEventAtStartDateResponse, err error) {
	query := `query ($sentiance_user_id: String!, $start_date: [String]!) {
		user(id: $sentiance_user_id) {
			id
			... on User {
				external_id
			}
			event_history(start: $start_date) {
				type
				start
				end
				... on Transport {
					event_id
					mode
					occupant_role
				}
			}
		}
	}`

	return res, c.Do(ctx, "POST", "/v2/gql", nil, GraphQLRequest{Query: query, Variables: req}, &res)
}
