package intercom

import (
	"encoding/json"
	"time"

	intercomSDK "gopkg.in/intercom/intercom-go.v2"
)

type Client struct {
	intercomSDK *intercomSDK.Client
}

// New creates a new Intercom Service, connecting with the access token provided.
func New(accessToken string) *Client {
	return &Client{
		intercomSDK: intercomSDK.NewClient(accessToken, ""),
	}
}

// SendInAppMessage sends a new message to a user from an admin. The body can contain both
// HTML and plain text information. The to field must be a Cuvva user id!
func (svc *Client) SendInAppMessage(to, from, body string) error {
	admin := intercomSDK.Admin{
		ID: json.Number(from),
	}
	user := intercomSDK.User{
		UserID: to,
	}

	msg := intercomSDK.NewInAppMessage(admin, user, body)
	_, err := svc.intercomSDK.Messages.Save(&msg)

	return err
}

// SendEvent accepts a UserID and EventName and sends a new event to Intercom
// It populates the event with a CreatedAt timestamp
func (svc *Client) SendEvent(userID, eventName string, metadata map[string]interface{}) error {
	return svc.intercomSDK.Events.Save(&intercomSDK.Event{
		UserID:    userID,
		EventName: eventName,
		CreatedAt: time.Now().Unix(),
		Metadata:  metadata,
	})
}

// SaveUser accepts a UserID and CustomAttributes and forwards the request to Intercom
// Intercom will create or update the user
func (svc *Client) SaveUser(userID string, customAttributes map[string]interface{}) error {
	_, err := svc.intercomSDK.Users.Save(&intercomSDK.User{
		UserID:           userID,
		CustomAttributes: customAttributes,
	})

	return err
}
