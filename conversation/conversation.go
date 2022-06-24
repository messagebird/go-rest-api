package conversation

import (
	"fmt"
	"net/http"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

const (
	// ConversationStatusActive is returned when the Conversation is active.
	// Only one active conversation can ever exist for a given contact.
	ConversationStatusActive ConversationStatus = "active"

	// ConversationStatusArchived is returned when the Conversation is
	// archived. When this is the case, a new Conversation is created when a
	// message is received from a contact.
	ConversationStatusArchived ConversationStatus = "archived"
)

type Conversation struct {
	ID                   string
	ContactID            string
	Contact              *Contact
	LastUsedChannelID    string
	Channels             []*Channel
	Messages             *MessagesCount
	Status               ConversationStatus
	CreatedDatetime      time.Time
	UpdatedDatetime      *time.Time
	LastReceivedDatetime *time.Time
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

type ConversationList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Conversation
}

// StartRequest contains the request data for the Start endpoint.
type StartRequest struct {
	ChannelID string                 `json:"channelId"`
	Content   *MessageContent        `json:"content"`
	To        MessageRecipient       `json:"to"`
	Type      MessageType            `json:"type"`
	Source    map[string]interface{} `json:"source,omitempty"`
	ReportUrl string                 `json:"reportUrl,omitempty"`
	Tag       MessageTag             `json:"tag,omitempty"`
	TrackId   string                 `json:"trackId,omitempty"`
	EventType string                 `json:"eventType,omitempty"`
	TTL       string                 `json:"ttl,omitempty"`
}

// ReplyRequest contains the request data for the Reply endpoint.
type ReplyRequest struct {
	Type      MessageType            `json:"type"`
	Content   *MessageContent        `json:"content"`
	ChannelID string                 `json:"channelId,omitempty"`
	Fallback  *Fallback              `json:"fallback,omitempty"`
	Source    map[string]interface{} `json:"source,omitempty"`
	EventType string                 `json:"eventType,omitempty"`
	ReportUrl string                 `json:"reportUrl,omitempty"`
	Tag       MessageTag             `json:"tag,omitempty"`
	TrackId   string                 `json:"trackId,omitempty"`
	TTL       string                 `json:"ttl,omitempty"`
}

type Fallback struct {
	From  string `json:"from"`
	After string `json:"after"`
}

// UpdateRequest contains the request data for the Update endpoint.
type UpdateRequest struct {
	Status ConversationStatus `json:"status"`
}

// DefaultListOptions provide a reasonable default for paginated endpoints.
var DefaultListOptions = &ListRequestOptions{10, 0}

// List gets a collection of Conversations. Pagination can be set in options.
func List(c *messagebird.Client, options *ListRequestOptions) (*ConversationList, error) {
	query := paginationQuery(options)

	convList := &ConversationList{}
	if err := request(c, convList, http.MethodGet, fmt.Sprintf("%s?%s", path, query), nil); err != nil {
		return nil, err
	}

	return convList, nil
}

// Read fetches a single Conversation based on its ID.
func Read(c *messagebird.Client, id string) (*Conversation, error) {
	conv := &Conversation{}
	if err := request(c, conv, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return conv, nil
}

// Start creates a conversation by sending an initial message. If an active
// conversation exists for the recipient, it is resumed.
func Start(c *messagebird.Client, req *StartRequest) (*Conversation, error) {
	conv := &Conversation{}
	if err := request(c, conv, http.MethodPost, path+"/start", req); err != nil {
		return nil, err
	}

	return conv, nil
}

// Reply Send a new message to an existing conversation. In case the conversation is archived, a new conversation is created.
func Reply(c *messagebird.Client, conversationId string, req *ReplyRequest) (*Message, error) {
	uri := fmt.Sprintf("%s/%s/%s", path, conversationId, messagesPath)

	message := &Message{}
	if err := request(c, message, http.MethodPost, uri, req); err != nil {
		return nil, err
	}

	return message, nil
}

// Update changes the conversation's status, so this can be used to (un)archive
// conversations.
func Update(c *messagebird.Client, id string, req *UpdateRequest) (*Conversation, error) {
	conv := &Conversation{}
	if err := request(c, conv, http.MethodPatch, path+"/"+id, req); err != nil {
		return nil, err
	}

	return conv, nil
}
