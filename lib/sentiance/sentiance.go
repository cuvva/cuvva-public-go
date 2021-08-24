package sentiance

import (
	"time"
)

type UserLink struct {
	ExternalID string `json:"external_id"`
}

type UserLinkResponse struct {
	ID string `json:"id"`
}

type UserLinkRequest struct {
	InstallID  string
	ExternalID string
}

type GetEventAtStartDateRequest struct {
	SentianceUserID string    `json:"sentiance_user_id"`
	StartDate       time.Time `json:"start_date"`
}

type GetEventByIDRequest struct {
	SentianceUserID string `json:"sentiance_user_id"`
	EventID         string `json:"event_id"`
}

type GetEventAndWaypointsByIDRequest struct {
	SentianceUserID string `json:"sentiance_user_id"`
	EventID         string `json:"event_id"`
}

type EnrichTransportEventResponse struct {
	EventID      string `json:"event_id"`
	Type         string `json:"type"`
	Mode         string `json:"mode"`
	OccupantRole string `json:"occupant_role"`
	Start        string `json:"start"`
	End          string `json:"end"`
}

type GraphQLRequest struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type TransportEventResponse struct {
	Data DataResponse `json:"data"`
}

type DataResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID           string                          `json:"id"`
	ExternalID   string                          `json:"external_id"`
	EventHistory []*EnrichTransportEventResponse `json:"event_history"`
}

type TransportEventWithWaypointsResponse struct {
	Data TransportEventWithWaypointsData `json:"data"`
}

type TransportEventWithWaypointsData struct {
	User UserWithWaypoints `json:"user"`
}

type UserWithWaypoints struct {
	ID           string             `json:"id"`
	ExternalID   string             `json:"external_id"`
	EventHistory []*TransportEvents `json:"event_history"`
}

type TransportEvents struct {
	EventID      string     `json:"event_id"`
	Type         string     `json:"type"`
	Mode         string     `json:"mode"`
	OccupantRole string     `json:"occupant_role"`
	AnalysisType string     `json:"analysis_type"`
	Start        string     `json:"start"`
	End          string     `json:"end"`
	Waypoints    []Waypoint `json:"waypoints"`
}

type Waypoint struct {
	Type      string    `json:"type"`
	Latitude  int64     `json:"longitude"`
	Longitude int64     `json:"latitude"`
	Timestamp time.Time `json:"timestamp"`
	Accuracy  int       `json:"accuracy"`
	Speed     int64     `json:"speed"`
	Altitude  int64     `json:"altitude"`
}
