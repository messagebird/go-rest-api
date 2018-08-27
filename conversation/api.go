package conversation

import (
	"encoding/json"
	"fmt"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

const (
	// apiRoot is the absolute URL of the Converstations API. All paths are
	// relative to apiRoot (e.g.
	// https://conversations.messagebird.com/v1/webhooks).
	apiRoot = "https://conversations.messagebird.com/v1"

	// path is the path for the Conversation resource, relative to apiRoot.
	path = "conversations"

	// messagesPath is the path for the Message resource, relative to apiRoot
	// and path.
	messagesPath = "messages"

	// webhooksPath is the path for the Webhook resource, relative to apiRoot.
	webhooksPath = "webhooks"
)

type ConversationList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Conversation
}

type Conversation struct {
	ID                   string
	ContactID            string
	Contact              *Contact
	LastUsedChannelID    string
	Channels             []*Channel
	Messages             *MessagesCount
	Status               ConversationStatus
	CreatedDatetime      *time.Time
	UpdatedDatetime      *time.Time
	LastReceivedDatetime *time.Time
}

type Contact struct {
	ID            string
	Href          string
	MSISDN        string
	FirstName     string
	LastName      string
	CustomDetails map[string]interface{}
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Channel struct {
	ID              string
	Name            string
	PlatformID      string
	Status          string
	CreatedDatetime *time.Time
	UpdatedDatetime *time.Time
}

type MessagesCount struct {
	HRef       string
	TotalCount int
}

// ConversationStatus indicates what state a Conversation is in.
type ConversationStatus string

const (
	// ConversationStatusActive is returned when the Conversation is active.
	// Only one active conversation can ever exist for a given contact.
	ConversationStatusActive ConversationStatus = "active"

	// ConversationStatusArchived is returned when the Conversation is
	// archived. When this is the case, a new Conversation is created when a
	// message is received from a contact.
	ConversationStatusArchived ConversationStatus = "archived"
)

type MessageList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Message
}

type Message struct {
	ID              string
	ConversationID  string
	ChannelID       string
	Direction       MessageDirection
	Status          MessageStatus
	Type            MessageType
	Content         MessageContent
	CreatedDatetime *time.Time
	UpdatedDatetime *time.Time
}

type MessageDirection string

const (
	// MessageDirectionReceived indicates an inbound message received from the customer.
	MessageDirectionReceived MessageDirection = "received"

	// MessageDirectionSent indicates an outbound message sent from the API.
	MessageDirectionSent MessageDirection = "sent"
)

// MessageStatus is a field set by the API. It indicates what the state of the
// message is, e.g. whether it has been successfully delivered or read.
type MessageStatus string

const (
	MessageStatusDeleted     MessageStatus = "deleted"
	MessageStatusDelivered   MessageStatus = "delivered"
	MessageStatusFailed      MessageStatus = "failed"
	MessageStatusPending     MessageStatus = "pending"
	MessageStatusRead        MessageStatus = "read"
	MessageStatusReceived    MessageStatus = "received"
	MessageStatusSent        MessageStatus = "sent"
	MessageStatusUnsupported MessageStatus = "unsupported"
)

// MessageType indicates what kind of content a Message has, e.g. audio or
// text.
type MessageType string

const (
	MessageTypeAudio    MessageType = "audio"
	MessageTypeFile     MessageType = "file"
	MessageTypeHSM      MessageType = "hsm"
	MessageTypeImage    MessageType = "image"
	MessageTypeLocation MessageType = "location"
	MessageTypeText     MessageType = "text"
	MessageTypeVideo    MessageType = "video"
)

// MessageContent holds a message's actual content. Only one field can be set
// per request.
type MessageContent struct {
	Audio    *Audio    `json:"audio,omitempty"`
	File     *File     `json:"file,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	Location *Location `json:"location,omitempty"`
	Video    *Video    `json:"video,omitempty"`
	Text     string    `json:"text,omitempty"`

	// HSM is a highly structured message for WhatsApp. Its definition lives in
	// hsm.go.
	HSM *HSM `json:"hsm,omitempty"`
}

type Media struct {
	URL string `json:"url"`
}

type Audio Media
type File Media
type Image Media
type Video Media

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type WebhookList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Webhook
}

type Webhook struct {
	ID              string
	ChannelID       string
	Events          []WebhookEvent
	URL             string
	CreatedDatetime *time.Time
	UpdatedDatetime *time.Time
}

type WebhookEvent string

const (
	WebhookEventConversationCreated WebhookEvent = "conversation.created"
	WebhookEventConversationUpdated WebhookEvent = "conversation.updated"
	WebhookEventMessageCreated      WebhookEvent = "message.created"
	WebhookEventMessageUpdated      WebhookEvent = "message.updated"
)

// request does the exact same thing as Client.Request. It does, however,
// prefix the path with the Conversation API's root. This ensures the client
// doesn't "handle" this for us: by default, it uses the REST API.
func request(c *messagebird.Client, v interface{}, method, path string, data interface{}) error {
	return c.Request(v, method, fmt.Sprintf("%s/%s", apiRoot, path), data)
}

// UnmarshalJSON is used to unmarshal the MSISDN to a string rather than an
// int64. The API returns integers, but this client always uses strings.
// Exposing a json.Number doesn't seem nice.
func (c *Contact) UnmarshalJSON(data []byte) error {
	target := struct {
		ID            string
		Href          string
		MSISDN        json.Number
		FirstName     string
		LastName      string
		CustomDetails map[string]interface{}
		CreatedAt     *time.Time
		UpdatedAt     *time.Time
	}{}

	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}

	// In many cases, the CustomDetails will contain the user ID. As
	// CustomDetails has interface{} values, these are unmarshalled as floats.
	// Convert them to int64.
	// Map key is not a typo: API returns userId and not userID.
	if val, ok := target.CustomDetails["userId"]; ok {
		var userID float64
		if userID, ok = val.(float64); ok {
			target.CustomDetails["userId"] = int64(userID)
		}
	}

	*c = Contact{
		target.ID,
		target.Href,
		target.MSISDN.String(),
		target.FirstName,
		target.LastName,
		target.CustomDetails,
		target.CreatedAt,
		target.UpdatedAt,
	}

	return nil
}
