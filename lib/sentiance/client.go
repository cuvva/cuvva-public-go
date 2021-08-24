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
	GetEventAtStartDate(ctx context.Context, req *GetEventAtStartDateRequest) (*TransportEventResponse, error)
	GetEventByID(ctx context.Context, req *GetEventByIDRequest) (*TransportEventResponse, error)
	GetEventAndWaypointsByID(ctx context.Context, req *GetEventByIDRequest) (*TransportEventWithWaypointsResponse, error)
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

func (c *HTTPClient) GetEventAtStartDate(ctx context.Context, req *GetEventAtStartDateRequest) (res *TransportEventResponse, err error) {
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

func (c *HTTPClient) GetEventByID(ctx context.Context, req *GetEventByIDRequest) (res *TransportEventResponse, err error) {
	query := `query ($sentiance_user_id: String!, $event_id: [String]!) {
		user(id: $sentiance_user_id) {
			id
			... on User {
				external_id
			}
			event_history(event_id: $event_id) {
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

func (c *HTTPClient) GetEventAndWaypointsByID(ctx context.Context, req *GetEventByIDRequest) (res *TransportEventWithWaypointsResponse, err error) {
	query := `query ($sentiance_user_id: String!, $event_id: [String]!) {
		user(id: $sentiance_user_id) {
			id
			... on User {
				external_id
			}
			event_history(event_id: $event_id) {
				type
				start
				end
				... on Transport {
					event_id
					mode
					occupant_role
					analysis_type
					waypoints {
						type
						latitude
						longitude
						timestamp
						accuracy
						speed
						altitude
					}
				}
			}
		}
	}`

	return res, c.Do(ctx, "POST", "/v2/gql", nil, GraphQLRequest{Query: query, Variables: req}, &res)
}
