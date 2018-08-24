package conversation

import (
	"encoding/json"
	"fmt"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

const (
	apiRoot      = "https://conversations.messagebird.com/v1"
	path         = "conversations"
	messagesPath = "messages"
	webhooksPath = "webhooks"
)

type ConversationList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      struct {
		First    string
		Previous string
		Next     string
		Last     string
	}
	Items []*Conversation
}

type Conversation struct {
	ID                   string
	ContactID            string
	Contact              *Contact
	LastUsedChannelID    string
	Channels             []*Channel
	Messages             *MessagesCount
	Status               ConversationStatus
	CreatedDatetime      time.Time
	UpdatedDatetime      time.Time
	LastReceivedDatetime time.Time
}

type Contact struct {
	ID            string
	Href          string
	MSISDN        string
	FirstName     string
	LastName      string
	CustomDetails map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Channel struct {
	ID              string
	Name            string
	PlatformID      string
	Status          string
	CreatedDatetime time.Time
	UpdatedDatetime time.Time
}

type MessagesCount struct {
	HRef       string
	TotalCount int
}

// ConversationStatus indicates what state a Conversation is in.
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"
	ConversationStatusArchived ConversationStatus = "archived"
)

type MessageList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      struct {
		First    string
		Previous string
		Next     string
		Last     string
	}
	Items []*Message
}

type Message struct {
	ID              string
	ConversationID  string
	ChannelID       string
	Direction       MessageDirection
	Status          MessageStatus
	Type            MessageType
	Content         MessageContent
	CreatedDatetime time.Time
	UpdatedDatetime time.Time
}

type WebhookList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      struct {
		First    string
		Previous string
		Next     string
		Last     string
	}
	Items []*Webhook
}

type Webhook struct {
	ID              string
	ChannelID       string
	Events          []WebhookEvent
	URL             string
	CreatedDatetime time.Time
	UpdatedDatetime time.Time
}

type MessageContent struct {
	Audio    *Audio    `json:"audio,omitempty"`
	HSM      *HSM      `json:"hsm,omitempty"`
	File     *File     `json:"file,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	Location *Location `json:"location,omitempty"`
	Video    *Video    `json:"video,omitempty"`
	Text     string    `json:"text,omitempty"`
}

type Media struct {
	URL string `json:"url"`
}

type Audio Media
type File Media
type Image Media
type Video Media

type HSM struct {
	Namespace             string                     `json:"namespace"`
	TemplateName          string                     `json:"templateName"`
	Language              *HSMLanguage               `json:"language"`
	LocalizableParameters []*HSMLocalizableParameter `json:"params"`
}

type HSMLanguage struct {
	Policy HSMLanguagePolicy `json:"policy"`
	Code   string            `json:"code"`
}

type HSMLanguagePolicy string

type HSMLocalizableParameter struct {
	Default  string                           `json:"default"`
	Currency *HSMLocalizableParameterCurrency `json:"currency,omitempty"`
	DateTime *time.Time                       `json:"dateTime,omitempty"`
}

type HSMLocalizableParameterCurrency struct {
	Code   string `json:"currencyCode"`
	Amount int64  `json:"amount"`
}

const (
	HSMLanguagePolicyFallback      HSMLanguagePolicy = "fallback"
	HSMLanguagePolicyDeterministic HSMLanguagePolicy = "deterministic"
)

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

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

type MessageDirection string

const (
	// MessageDirectionReceived indicates an inbound message received from the customer.
	MessageDirectionReceived MessageDirection = "received"

	// MessageDirectionSent indicates an outbound message sent from the API.
	MessageDirectionSent MessageDirection = "sent"
)

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

type WebhookEvent string

const (
	WebhookEventConversationCreated WebhookEvent = "conversation.created"
	WebhookEventConversationUpdated WebhookEvent = "conversation.updated"
	WebhookEventMessageCreated      WebhookEvent = "message.created"
	WebhookEventMessageUpdated      WebhookEvent = "message.updated"
)

// request does the exact same thign as Client.Request. It does, however,
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
		CreatedAt     time.Time
		UpdatedAt     time.Time
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
