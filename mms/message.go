package mms

import (
	"net/http"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

// path represents the path to the MMS resource.
const path = "mms"

// Message represents a MMS Message.
type Message struct {
	ID                string
	HRef              string
	Direction         string
	Originator        string
	Body              string
	Reference         string
	Subject           string
	MediaUrls         []string
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        messagebird.Recipients
}

type CreateRequest struct {
	Originator        string     `json:"originator"` // the sender of the message.
	Recipients        string     `json:"recipients"` // comma separated list
	Body              string     `json:"body,omitempty"`
	MediaUrls         []string   `json:"mediaUrls"`
	Subject           string     `json:"subject,omitempty"`
	Reference         string     `json:"reference,omitempty"`
	ScheduledDatetime *time.Time `json:"scheduledDatetime,omitempty"`
}

// Read retrieves the information of an existing MmsMessage.
func Read(c *messagebird.Client, id string) (*Message, error) {
	mmsMessage := &Message{}
	if err := c.Request(mmsMessage, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return mmsMessage, nil
}

// Create creates a new MMS message for one or more recipients.
// Max of 50 recipients can be entered per request.
func Create(c *messagebird.Client, req *CreateRequest) (*Message, error) {
	mmsMessage := &Message{}
	if err := c.Request(mmsMessage, http.MethodPost, path, req); err != nil {
		return nil, err
	}

	return mmsMessage, nil
}
